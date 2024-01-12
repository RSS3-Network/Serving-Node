package multicall3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

// Multicall https://github.com/mds1/multicall
// https://etherscan.io/address/0xcA11bde05977b3631167028862bE2a173976CA11
//go:generate go run -mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Multicall3.abi --pkg multicall3 --type Multicall3 --out contract_multicall3.go

var (
	AddressMulticall3 = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
)

var deployedAtMap = map[uint64]uint64{
	uint64(filter.EthereumChainIDMainnet):  14353601, // https://etherscan.io/tx/0x00d9fcb7848f6f6b0aae4fb709c133d69262b902156c85a473ef23faa60760bd
	uint64(filter.EthereumChainIDOptimism): 4286263,  // https://optimistic.etherscan.io/tx/0xb62f9191a2cf399c0d2afd33f5b8baf7c6b52af6dd2386e44121b1bab91b80e5
	uint64(filter.EthereumChainIDPolygon):  15921452, // https://bscscan.com/tx/0xcc0ddf5f791617ba9befce57995dbcb3a202946a1eefa3469742b01a0decdaf2
}

func IsDeployed(chainID uint64, blockNumber *big.Int) bool {
	deployedAt, exists := deployedAtMap[chainID]
	if !exists {
		return false
	}

	if blockNumber == nil {
		return true
	}

	return deployedAt < blockNumber.Uint64()
}

func Aggregate3(ctx context.Context, chainID uint64, calls []Multicall3Call3, blockNumber *big.Int, contractBackend bind.ContractCaller) ([]*Multicall3Result, error) {
	// If the Multicall3 contract is not yet deployed, split it into multiple regular requests.
	if !IsDeployed(chainID, blockNumber) {
		resultPool := pool.NewWithResults[*Multicall3Result]().WithContext(ctx).WithCancelOnError()

		for _, call := range calls {
			call := call

			resultPool.Go(func(ctx context.Context) (*Multicall3Result, error) {
				message := ethereum.CallMsg{
					To:   lo.ToPtr(call.Target),
					Data: call.CallData,
				}

				data, err := contractBackend.CallContract(ctx, message, blockNumber)
				if err != nil && !call.AllowFailure {
					return nil, err
				}

				result := Multicall3Result{
					Success:    err == nil,
					ReturnData: data,
				}

				return &result, nil
			})
		}

		return resultPool.Wait()
	}

	abi, err := Multicall3MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("load abi: %w", err)
	}

	callData, err := abi.Pack("aggregate3", calls)
	if err != nil {
		return nil, fmt.Errorf("pack data: %w", err)
	}

	message := ethereum.CallMsg{
		To:   &AddressMulticall3,
		Data: callData,
	}

	results := make([]Multicall3Result, 0, len(calls))

	data, err := contractBackend.CallContract(ctx, message, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("call contract: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("data in empty")
	}

	if err := abi.UnpackIntoInterface(&results, "aggregate3", data); err != nil {
		return nil, fmt.Errorf("unpack result: %w", err)
	}

	return lo.ToSlicePtr(results), nil
}
