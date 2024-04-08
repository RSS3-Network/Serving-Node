package curve

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/rueidis"
	"github.com/rss3-network/node/config"
	"github.com/rss3-network/node/internal/engine"
	source "github.com/rss3-network/node/internal/engine/source/ethereum"
	"github.com/rss3-network/node/internal/engine/worker/contract/curve/pool"
	"github.com/rss3-network/node/provider/ethereum"
	"github.com/rss3-network/node/provider/ethereum/contract"
	"github.com/rss3-network/node/provider/ethereum/contract/curve"
	"github.com/rss3-network/node/provider/ethereum/contract/erc20"
	"github.com/rss3-network/node/provider/ethereum/token"
	"github.com/rss3-network/node/provider/httpx"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Worker is the worker for Curve.
var _ engine.Worker = (*worker)(nil)

type worker struct {
	config                        *config.Module
	ethereumClient                ethereum.Client
	tokenClient                   token.Client
	curvePoolRegistry             pool.Registry
	curveRegistryExchangeFilterer *curve.RegistryExchangeFilterer
	curveStableSwapFilterer       *curve.StableSwapFilterer
	curveStableSwapNCoinsFilterer *curve.StableSwapNCoinsFilterer
	curveLiquidityGaugeFilterer   *curve.LiquidityGaugeFilterer
	erc20ERC20Filterer            *erc20.ERC20Filterer
}

func (w *worker) Name() string {
	return filter.Curve.String()
}

// Filter curve contract address and event hash.
func (w *worker) Filter() engine.SourceFilter {
	return &source.Filter{
		LogAddresses: []common.Address{
			curve.AddressRegistryExchangeEthereum,
			curve.AddressRegistryExchangeArbitrum,
			curve.AddressRegistryExchangeAvalanche,
			curve.AddressRegistryExchangeGnosis,
			curve.AddressRegistryExchangeOptimism,
			curve.AddressRegistryExchangePolygon,
		},
		LogTopics: []common.Hash{
			curve.EventHashStableSwapTokenExchange,
			curve.EventHashStableSwapAddLiquidity2Coins,
			curve.EventHashStableSwapAddLiquidity3Coins,
			curve.EventHashStableSwapAddLiquidity4Coins,
			curve.EventHashStableSwapRemoveLiquidity2Coins,
			curve.EventHashStableSwapRemoveLiquidity3Coins,
			curve.EventHashStableSwapRemoveLiquidity4Coins,
			curve.EventHashStableSwapRemoveLiquidityOne,
			curve.EventHashStableSwapRemoveLiquidityOneFactory,
			curve.EventHashStableSwapRemoveLiquidityImbalance2Coins,
			curve.EventHashStableSwapRemoveLiquidityImbalance3Coins,
			curve.EventHashStableSwapRemoveLiquidityImbalance4Coins,
			curve.EventHashLiquidityGaugeDeposit,
			curve.EventHashLiquidityGaugeWithdraw,
		},
	}
}

func (w *worker) Match(_ context.Context, task engine.Task) (bool, error) {
	return task.GetNetwork().Source() == filter.NetworkEthereumSource, nil
}

// Transform Ethereum task to feed.
func (w *worker) Transform(ctx context.Context, task engine.Task) (*schema.Feed, error) {
	ethereumTask, ok := task.(*source.Task)
	if !ok {
		return nil, fmt.Errorf("invalid task type: %T", task)
	}

	// Build default curve feed from task.
	feed, err := ethereumTask.BuildFeed(schema.WithFeedPlatform(filter.PlatformCurve))
	if err != nil {
		return nil, fmt.Errorf("build feed: %w", err)
	}

	switch {
	case w.matchStableSwapEventHashStableSwapAddLiquidityTransaction(ethereumTask):
		feed.Type = filter.TypeExchangeLiquidity

		actions, err := w.transformStableSwapAddLiquidityTransaction(ctx, ethereumTask)
		if err != nil {
			return nil, fmt.Errorf("handle ethereum stable swap add liquidity transaction: %w", err)
		}

		feed.Actions = append(feed.Actions, actions...)

		return feed, nil
	case w.matchStableSwapEventHashStableSwapRemoveLiquidityTransaction(ethereumTask):
		feed.Type = filter.TypeExchangeLiquidity

		actions, err := w.transformStableSwapRemoveLiquidityTransaction(ctx, ethereumTask)
		if err != nil {
			return nil, fmt.Errorf("handle ethereum stable swap remove liquidity transaction: %w", err)
		}

		feed.Actions = append(feed.Actions, actions...)

		return feed, nil
	}

	// Match and handle ethereum logs.
	for _, log := range ethereumTask.Receipt.Logs {
		// Ignore anonymous logs.
		if len(log.Topics) == 0 {
			continue
		}

		var (
			actions []*schema.Action
			err     error
		)

		// Match curve core contract events
		switch {
		case w.matchEthereumRegistryExchangeExchangeMultipleLog(ethereumTask, log):
			// Add ETH liquidity
			feed.Type = filter.TypeExchangeSwap
			actions, err = w.transformRegistryExchangeExchangeMultipleLog(ctx, ethereumTask, log)
		case w.matchEthereumStableSwapTokenExchangeLog(ethereumTask, log):
			// Add stable swap token exchange
			feed.Type = filter.TypeExchangeSwap
			actions, err = w.transformStableSwapTokenExchangeLog(ctx, ethereumTask, log)
		case w.matchEthereumLiquidityGaugeDepositLog(ethereumTask, log):
			// Add liquidity gauge deposit
			feed.Type = filter.TypeExchangeStaking
			actions, err = w.transformLiquidityGaugeDepositLog(ctx, ethereumTask, log)
		case w.matchEthereumLiquidityGaugeWithdrawLog(ethereumTask, log):
			// Add liquidity gauge withdraw
			feed.Type = filter.TypeExchangeStaking
			actions, err = w.transformLiquidityGaugeWithdrawLog(ctx, ethereumTask, log)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		feed.Actions = append(feed.Actions, actions...)
	}

	return feed, nil
}

func (w *worker) matchStableSwapEventHashStableSwapAddLiquidityTransaction(task *source.Task) bool {
	return contract.MatchMethodIDs(
		task.Transaction.Input,
		curve.MethodIDStableSwapAddLiquidity2Coins,
		curve.MethodIDStableSwapAddLiquidity3Coins,
		curve.MethodIDStableSwapAddLiquidity4Coins,
	)
}

func (w *worker) matchStableSwapEventHashStableSwapRemoveLiquidityTransaction(task *source.Task) bool {
	return contract.MatchMethodIDs(
		task.Transaction.Input,
		curve.MethodIDStableSwapRemoveLiquidity2Coins,
		curve.MethodIDStableSwapRemoveLiquidity3Coins,
		curve.MethodIDStableSwapRemoveLiquidity4Coins,
		curve.MethodIDStableSwapRemoveLiquidityOne,
		curve.MethodIDStableSwapRemoveLiquidityImbalance2Coins,
		curve.MethodIDStableSwapRemoveLiquidityImbalance3Coins,
		curve.MethodIDStableSwapRemoveLiquidityImbalance4Coins,
	)
}

func (w *worker) matchEthereumRegistryExchangeExchangeMultipleLog(_ *source.Task, log *ethereum.Log) bool {
	return len(log.Topics) == 3 && contract.MatchEventHashes(log.Topics[0], curve.EventHashRegistryExchangeExchangeMultiple)
}

func (w *worker) matchEthereumStableSwapTokenExchangeLog(task *source.Task, log *ethereum.Log) bool {
	if len(log.Topics) != 2 || !contract.MatchEventHashes(log.Topics[0], curve.EventHashStableSwapTokenExchange) {
		return false
	}

	validated, err := w.curvePoolRegistry.Validate(context.Background(), task.Network, pool.ContractTypePool, log.Address)
	if err != nil {
		zap.L().Error("validate pool", zap.Error(err), zap.Stringer("pool", log.Address))
	}

	return validated != nil
}

func (w *worker) matchEthereumLiquidityGaugeDepositLog(task *source.Task, log *ethereum.Log) bool {
	if len(log.Topics) != 2 || !contract.MatchEventHashes(log.Topics[0], curve.EventHashLiquidityGaugeDeposit) {
		return false
	}

	validated, err := w.curvePoolRegistry.Validate(context.Background(), task.Network, pool.ContractTypeGauge, log.Address)
	if err != nil {
		zap.L().Error("validate gauge", zap.Error(err), zap.Stringer("pool", log.Address))
	}

	return validated != nil
}

func (w *worker) matchEthereumLiquidityGaugeWithdrawLog(task *source.Task, log *ethereum.Log) bool {
	if len(log.Topics) != 2 || !contract.MatchEventHashes(log.Topics[0], curve.EventHashLiquidityGaugeWithdraw) {
		return false
	}

	validated, err := w.curvePoolRegistry.Validate(context.Background(), task.Network, pool.ContractTypeGauge, log.Address)
	if err != nil {
		zap.L().Error("validate gauge", zap.Error(err), zap.Stringer("pool", log.Address))
	}

	return validated != nil
}

func (w *worker) transformStableSwapAddLiquidityTransaction(ctx context.Context, task *source.Task) ([]*schema.Action, error) {
	addLiquidityLog, _ := lo.Find(task.Receipt.Logs, func(log *ethereum.Log) bool {
		return len(log.Topics) > 0 && contract.MatchEventHashes(
			log.Topics[0],
			curve.EventHashStableSwapAddLiquidity2Coins,
			curve.EventHashStableSwapAddLiquidity3Coins,
			curve.EventHashStableSwapAddLiquidity4Coins,
		)
	})

	var (
		poolAddress = addLiquidityLog.Address
		actions     = make([]*schema.Action, 0)
	)

	for _, log := range task.Receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		if !contract.MatchEventHashes(log.Topics[0], erc20.EventHashTransfer) {
			continue
		}

		transferEvent, err := w.erc20ERC20Filterer.ParseTransfer(log.Export())
		if err != nil {
			return nil, fmt.Errorf("parse transfer event: %w", err)
		}

		switch {
		case transferEvent.From == task.Transaction.From && transferEvent.To == poolAddress: // Add liquidity
			action, err := w.buildEthereumExchangeLiquidityAction(ctx, task.Header.Number, task.ChainID, transferEvent.From, transferEvent.To, &transferEvent.Raw.Address, transferEvent.Value, metadata.ActionExchangeLiquidityAdd)
			if err != nil {
				return nil, fmt.Errorf("build exchange liquidity action: %w", err)
			}

			actions = append(actions, action)
		case transferEvent.From == ethereum.AddressGenesis && transferEvent.To == task.Transaction.From: // Mint coin
			action, err := w.buildEthereumTransactionTransferAction(ctx, task.Header.Number, task.ChainID, transferEvent.From, transferEvent.To, &transferEvent.Raw.Address, transferEvent.Value)
			if err != nil {
				return nil, fmt.Errorf("build transaction transfer action: %w", err)
			}

			actions = append(actions, action)
		}
	}

	return actions, nil
}

func (w *worker) transformStableSwapRemoveLiquidityTransaction(ctx context.Context, task *source.Task) ([]*schema.Action, error) {
	addLiquidityLog, _ := lo.Find(task.Receipt.Logs, func(log *ethereum.Log) bool {
		return len(log.Topics) > 0 && contract.MatchEventHashes(
			log.Topics[0],
			curve.EventHashStableSwapRemoveLiquidity2Coins,
			curve.EventHashStableSwapRemoveLiquidity3Coins,
			curve.EventHashStableSwapRemoveLiquidity4Coins,
			curve.EventHashStableSwapRemoveLiquidityOne,
			curve.EventHashStableSwapRemoveLiquidityOneFactory,
			curve.EventHashStableSwapRemoveLiquidityImbalance2Coins,
			curve.EventHashStableSwapRemoveLiquidityImbalance3Coins,
			curve.EventHashStableSwapRemoveLiquidityImbalance4Coins,
		)
	})

	var (
		poolAddress = addLiquidityLog.Address
		actions     = make([]*schema.Action, 0)
	)

	for _, log := range task.Receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		if !contract.MatchEventHashes(log.Topics[0], erc20.EventHashTransfer) {
			continue
		}

		transferEvent, err := w.erc20ERC20Filterer.ParseTransfer(log.Export())
		if err != nil {
			return nil, fmt.Errorf("parse transfer event: %w", err)
		}

		switch {
		case transferEvent.From == poolAddress && transferEvent.To == task.Transaction.From: // Remove liquidity
			action, err := w.buildEthereumExchangeLiquidityAction(ctx, task.Header.Number, task.ChainID, transferEvent.From, transferEvent.To, &transferEvent.Raw.Address, transferEvent.Value, metadata.ActionExchangeLiquidityRemove)
			if err != nil {
				return nil, fmt.Errorf("build exchange liquidity action: %w", err)
			}

			actions = append(actions, action)
		case transferEvent.From == task.Transaction.From && transferEvent.To == ethereum.AddressGenesis: // Burn coin
			action, err := w.buildEthereumTransactionTransferAction(ctx, task.Header.Number, task.ChainID, transferEvent.From, transferEvent.To, &transferEvent.Raw.Address, transferEvent.Value)
			if err != nil {
				return nil, fmt.Errorf("build transaction transfer action: %w", err)
			}

			actions = append(actions, action)
		}
	}

	return actions, nil
}

func (w *worker) transformRegistryExchangeExchangeMultipleLog(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	event, err := w.curveRegistryExchangeFilterer.ParseExchangeMultiple(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse exchange multiple event: %w", err)
	}

	var tokenIn, tokenOut common.Address

	tokenIn = event.Route[0]

	for index := len(event.Route) - 1; index >= 0; index-- {
		if router := event.Route[index]; router != ethereum.AddressGenesis {
			tokenOut = event.Route[index]

			break
		}
	}

	action, err := w.buildEthereumExchangeSwapAction(ctx, task.Header.Number, task.ChainID, event.Buyer, event.Receiver, tokenIn, tokenOut, event.AmountSold, event.AmountBought)
	if err != nil {
		return nil, fmt.Errorf("build exchange swap action: %w", err)
	}

	actions := []*schema.Action{
		action,
	}

	return actions, nil
}

func (w *worker) transformStableSwapTokenExchangeLog(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	event, err := w.curveStableSwapFilterer.ParseTokenExchange(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse token exchange event: %w", err)
	}

	stableSwapCaller, err := curve.NewStableSwapCaller(log.Address, w.ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new StableSwap caller: %w", err)
	}

	tokenIn, err := stableSwapCaller.Coins(nil, event.SoldId)
	if err != nil {
		return nil, fmt.Errorf("get token by sold id: %w", err)
	}

	tokenOut, err := stableSwapCaller.Coins(nil, event.BoughtId)
	if err != nil {
		return nil, fmt.Errorf("get token by bought id: %w", err)
	}

	action, err := w.buildEthereumExchangeSwapAction(ctx, task.Header.Number, task.ChainID, event.Buyer, event.Buyer, tokenIn, tokenOut, event.TokensSold, event.TokensBought)
	if err != nil {
		return nil, fmt.Errorf("build exchange swap action: %w", err)
	}

	actions := []*schema.Action{
		action,
	}

	return actions, nil
}

func (w *worker) transformLiquidityGaugeDepositLog(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	event, err := w.curveLiquidityGaugeFilterer.ParseDeposit(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse deposit event: %w", err)
	}

	liquidityGaugeCaller, err := curve.NewLiquidityGaugeCaller(log.Address, w.ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new LiquidityGauge caller: %w", err)
	}

	callOptions := bind.CallOpts{
		BlockNumber: task.Header.Number,
	}

	liquidityProviderTokenAddress, err := liquidityGaugeCaller.LpToken(&callOptions)
	if err != nil {
		return nil, fmt.Errorf("get lp token: %w", err)
	}

	action, err := w.buildEthereumTransactionStakingAction(ctx, task, event.Provider, event.Raw.Address, liquidityProviderTokenAddress, event.Value, metadata.ActionExchangeStakingStake, nil)
	if err != nil {
		return nil, fmt.Errorf("build exchange swap action: %w", err)
	}

	actions := []*schema.Action{
		action,
	}

	return actions, nil
}

func (w *worker) transformLiquidityGaugeWithdrawLog(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	event, err := w.curveLiquidityGaugeFilterer.ParseWithdraw(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse withdraw event: %w", err)
	}

	liquidityGaugeCaller, err := curve.NewLiquidityGaugeCaller(log.Address, w.ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new LiquidityGauge caller: %w", err)
	}

	callOptions := bind.CallOpts{
		BlockNumber: task.Header.Number,
	}

	liquidityProviderTokenAddress, err := liquidityGaugeCaller.LpToken(&callOptions)
	if err != nil {
		return nil, fmt.Errorf("get lp token: %w", err)
	}

	action, err := w.buildEthereumTransactionStakingAction(ctx, task, event.Raw.Address, event.Provider, liquidityProviderTokenAddress, event.Value, metadata.ActionExchangeStakingUnstake, nil)
	if err != nil {
		return nil, fmt.Errorf("build exchange swap action: %w", err)
	}

	actions := []*schema.Action{
		action,
	}

	return actions, nil
}

func (w *worker) buildEthereumExchangeLiquidityAction(ctx context.Context, blockNumber *big.Int, chainID uint64, sender, receiver common.Address, tokenAddress *common.Address, tokenValue *big.Int, liquidityAction metadata.ExchangeLiquidityAction) (*schema.Action, error) {
	tokenMetadata, err := w.tokenClient.Lookup(ctx, chainID, tokenAddress, nil, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup token: %w", err)
	}

	tokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(tokenValue, 0))

	action := schema.Action{
		Type:     filter.TypeExchangeLiquidity,
		Platform: filter.PlatformCurve.String(),
		From:     sender.String(),
		To:       receiver.String(),
		Metadata: metadata.ExchangeLiquidity{
			Action: liquidityAction,
			Tokens: []metadata.Token{
				*tokenMetadata,
			},
		},
	}

	return &action, nil
}

func (w *worker) buildEthereumTransactionTransferAction(ctx context.Context, blockNumber *big.Int, chainID uint64, sender, receiver common.Address, tokenAddress *common.Address, tokenValue *big.Int) (*schema.Action, error) {
	tokenMetadata, err := w.tokenClient.Lookup(ctx, chainID, tokenAddress, nil, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup token: %w", err)
	}

	tokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(tokenValue, 0))

	actionType := filter.TypeTransactionTransfer

	if sender == ethereum.AddressGenesis {
		actionType = filter.TypeTransactionMint
	}

	if receiver == ethereum.AddressGenesis {
		actionType = filter.TypeTransactionBurn
	}

	action := schema.Action{
		Type:     actionType,
		Platform: filter.PlatformCurve.String(),
		From:     sender.String(),
		To:       receiver.String(),
		Metadata: metadata.TransactionTransfer(*tokenMetadata),
	}

	return &action, nil
}

func (w *worker) buildEthereumExchangeSwapAction(ctx context.Context, blockNumber *big.Int, chainID uint64, from, to, tokenIn, tokenOut common.Address, amountIn, amountOut *big.Int) (*schema.Action, error) {
	tokenInMetadata, err := w.tokenClient.Lookup(ctx, chainID, &tokenIn, nil, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup token metadata %s: %w", tokenIn, err)
	}

	tokenInMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(amountIn, 0).Abs())

	tokenOutMetadata, err := w.tokenClient.Lookup(ctx, chainID, &tokenOut, nil, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup token metadata %s: %w", tokenOut, err)
	}

	tokenOutMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(amountOut, 0).Abs())

	action := schema.Action{
		Type:     filter.TypeExchangeSwap,
		Platform: filter.PlatformCurve.String(),
		From:     from.String(),
		To:       to.String(),
		Metadata: metadata.ExchangeSwap{
			From: *tokenInMetadata,
			To:   *tokenOutMetadata,
		},
	}

	return &action, nil
}

func (w *worker) buildEthereumTransactionStakingAction(ctx context.Context, task *source.Task, from, to common.Address, token common.Address, value *big.Int, stakingAction metadata.ExchangeStakingAction, period *metadata.ExchangeStakingPeriod) (*schema.Action, error) {
	tokenMetadata, err := w.tokenClient.Lookup(ctx, task.ChainID, &token, nil, task.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("lookup token: %w", err)
	}

	tokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(value, 0))

	action := schema.Action{
		Type:     filter.TypeExchangeStaking,
		Platform: filter.PlatformCurve.String(),
		From:     from.String(),
		To:       to.String(),
		Metadata: metadata.ExchangeStaking{
			Action: stakingAction,
			Token:  *tokenMetadata,
			Period: period,
		},
	}

	return &action, nil
}

// NewWorker creates a new Curve worker.
func NewWorker(config *config.Module, redisClient rueidis.Client) (engine.Worker, error) {
	var (
		err      error
		instance = worker{
			config: config,
		}
	)

	// Initialize ethereum client.
	if instance.ethereumClient, err = ethereum.Dial(context.Background(), config.Endpoint); err != nil {
		return nil, fmt.Errorf("initialize ethereum client: %w", err)
	}

	// Initialize token client.
	instance.tokenClient = token.NewClient(instance.ethereumClient)

	// Initialize curve filterers.
	instance.curveRegistryExchangeFilterer = lo.Must(curve.NewRegistryExchangeFilterer(ethereum.AddressGenesis, nil))
	instance.curveStableSwapFilterer = lo.Must(curve.NewStableSwapFilterer(ethereum.AddressGenesis, nil))
	instance.curveStableSwapNCoinsFilterer = lo.Must(curve.NewStableSwapNCoinsFilterer(ethereum.AddressGenesis, nil))
	instance.curveLiquidityGaugeFilterer = lo.Must(curve.NewLiquidityGaugeFilterer(ethereum.AddressGenesis, nil))
	instance.erc20ERC20Filterer = lo.Must(erc20.NewERC20Filterer(ethereum.AddressGenesis, nil))

	// Initialize curve registry.
	httpClient, err := httpx.NewHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("new http client: %w", err)
	}

	instance.curvePoolRegistry = pool.NewRegistry(redisClient, httpClient)

	return &instance, nil
}
