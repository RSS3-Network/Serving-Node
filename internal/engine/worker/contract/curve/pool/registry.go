package pool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/rueidis"
	"github.com/rss3-network/node/provider/httpx"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

type Registry interface {
	Refresh(ctx context.Context) error
	Validate(ctx context.Context, chain filter.Network, contractType ContractType, address common.Address) (*Pool, error)
}

var _ Registry = (*registry)(nil)

type registry struct {
	redisClient rueidis.Client
	httpClient  httpx.Client
}

func (r *registry) Refresh(ctx context.Context) error {
	const Endpoint = "https://api.curve.fi/"

	networks := []filter.Network{
		filter.NetworkArbitrum,
		filter.NetworkAvalanche,
		filter.NetworkEthereum,
		filter.NetworkGnosis,
		filter.NetworkOptimism,
		filter.NetworkPolygon,
	}

	registryIDs := []string{
		"main",
		"crypto",
		"factory",
		"factory-crypto",
	}

	resultPool := pool.NewWithResults[[]Pool]().
		WithContext(ctx).
		WithFirstError().
		WithCancelOnError()

	for _, network := range networks {
		network := network

		for _, registryID := range registryIDs {
			resultPool.Go(func(ctx context.Context) ([]Pool, error) {
				readCloser, err := r.httpClient.Fetch(ctx, fmt.Sprintf("%sapi/getPools/%s/%s", Endpoint, network.String(), registryID))
				if err != nil {
					return nil, fmt.Errorf("fetch request: %w", err)
				}
				defer lo.Try(readCloser.Close)

				var result Response[GetPoolData]

				if err := json.NewDecoder(readCloser).Decode(&result); err != nil {
					return nil, fmt.Errorf("decode json: %w", err)
				}

				for index, curvePool := range result.Data.PoolData {
					curvePool.Network = network

					result.Data.PoolData[index] = curvePool
				}

				return result.Data.PoolData, nil
			})
		}
	}

	result, err := resultPool.Wait()
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	curvePools := lo.Flatten(result)

	commands := make([]rueidis.Completed, 0, len(curvePools))

	for _, curvePool := range curvePools {
		keys := []string{
			r.formatRedisKey(curvePool.Network, ContractTypePool, curvePool.Address),
			r.formatRedisKey(curvePool.Network, ContractTypeToken, curvePool.LiquidityProviderTokenAddress),
			r.formatRedisKey(curvePool.Network, ContractTypeGauge, curvePool.GaugeAddress),
		}

		for _, key := range keys {
			command := r.redisClient.B().Set().Key(key).Value(curvePool.Name).Build()

			commands = append(commands, command)
		}
	}

	for _, redisResult := range r.redisClient.DoMulti(ctx, commands...) {
		if err := redisResult.Error(); err != nil {
			return fmt.Errorf("redis result: %w", err)
		}
	}

	return nil
}

func (r *registry) Validate(ctx context.Context, network filter.Network, contractType ContractType, address common.Address) (*Pool, error) {
	command := r.redisClient.B().Get().Key(r.formatRedisKey(network, contractType, address)).Build()

	result := r.redisClient.Do(ctx, command)

	if err := result.Error(); err != nil {
		return nil, fmt.Errorf("redis result: %w", err)
	}

	curvePool := Pool{
		Network: network,
		Address: address,
	}

	var err error
	if curvePool.Name, err = result.ToString(); err != nil {
		return nil, fmt.Errorf("to string: %w", err)
	}

	return &curvePool, nil
}

func (r *registry) formatRedisKey(network filter.Network, contractType ContractType, address common.Address) string {
	return fmt.Sprintf("curve:%s:%s:%s", network.String(), contractType, address)
}

func NewRegistry(redisClient rueidis.Client, httpClient httpx.Client) Registry {
	return &registry{
		redisClient: redisClient,
		httpClient:  httpClient,
	}
}
