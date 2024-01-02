package aavegotchi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/internal/engine"
	source "github.com/naturalselectionlabs/rss3-node/internal/engine/source/ethereum"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/aavegotchi"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc20"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/token"
	"github.com/naturalselectionlabs/rss3-node/schema"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/naturalselectionlabs/rss3-node/schema/metadata"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var _ engine.Worker = (*worker)(nil)

type worker struct {
	tokenClient                token.Client
	erc20Filterer              *erc20.ERC20Filterer
	erc721MarketplaceFilterer  *aavegotchi.ERC721MarketplaceFilterer
	erc1155MarketplaceFilterer *aavegotchi.ERC1155MarketplaceFilterer
}

func (w *worker) Name() string {
	return engine.Aavegotchi.String()
}

// Filter filters the source for Aavegotchi.
func (w *worker) Filter() engine.SourceFilter {
	return &source.Filter{
		LogAddresses: []common.Address{
			aavegotchi.AddressAavegotchi,
		},
		LogTopics: []common.Hash{
			aavegotchi.EventHashERC721ListingAdd,
			aavegotchi.EventHashERC721ExecutedListing,
			aavegotchi.EventHashERC1155ListingAdd,
			aavegotchi.EventHashERC1155ExecutedListing,
			erc20.EventHashTransfer,
		},
	}
}

func (w *worker) Match(_ context.Context, _ engine.Task) (bool, error) {
	// TODO Delete
	return true, nil
}

func (w *worker) Transform(ctx context.Context, task engine.Task) (*schema.Feed, error) {
	// Cast the task to an Ethereum task.
	ethereumTask, ok := task.(*source.Task)
	if !ok {
		return nil, fmt.Errorf("invalid task type: %T", task)
	}

	// Build the feed.
	feed, err := task.BuildFeed(schema.WithFeedPlatform(filter.PlatformAavegotchi))
	if err != nil {
		return nil, fmt.Errorf("build feed: %w", err)
	}

	for _, log := range ethereumTask.Receipt.Logs {
		var (
			action *schema.Action
			err    error
		)

		if len(log.Topics) == 0 {
			continue
		}

		switch {
		case w.matchERC721ListingAdd(ctx, *log):
			action, err = w.handleERC721ListingAdd(ctx, ethereumTask, *log, feed)
		case w.matchERC721ExecutedListing(ctx, *log):
			action, err = w.handleERC721ExecutedListing(ctx, ethereumTask, *log, feed)
		case w.matchERC1155ListingAdd(ctx, *log):
			action, err = w.handleERC1155ListingAdd(ctx, ethereumTask, *log, feed)
		case w.matchERC1155ExecutedListing(ctx, *log):
			action, err = w.handleERC1155ExecutedListing(ctx, ethereumTask, *log, feed)
		case w.matchERC20TransferLog(ctx, *log):
			action, err = w.handleERC20TransferLog(ctx, ethereumTask, *log, feed)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		feed.Actions = append(feed.Actions, action)
	}

	if feed.Type == filter.TypeMetaverseTrade {
		return w.handleMetaverseTradeCost(ctx, feed)
	}

	return feed, nil
}

func (w *worker) matchERC721ListingAdd(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == aavegotchi.EventHashERC721ListingAdd
}

func (w *worker) matchERC721ExecutedListing(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == aavegotchi.EventHashERC721ExecutedListing
}

func (w *worker) matchERC1155ListingAdd(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == aavegotchi.EventHashERC1155ListingAdd
}

func (w *worker) matchERC1155ExecutedListing(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == aavegotchi.EventHashERC1155ExecutedListing
}

func (w *worker) matchERC20TransferLog(_ context.Context, log ethereum.Log) bool {
	return len(log.Topics) == 3 && log.Topics[0] == erc20.EventHashTransfer
}

func (w *worker) handleERC721ListingAdd(ctx context.Context, task *source.Task, log ethereum.Log, feed *schema.Feed) (*schema.Action, error) {
	event, err := w.erc721MarketplaceFilterer.ParseERC721ListingAdd(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse erc721 listing add: %w", err)
	}

	feed.Type = filter.TypeMetaverseTrade

	return w.buildTradeAction(ctx, log.BlockNumber, task.ChainID, feed.Type, metadata.ActionMetaverseTradeList, event.Seller, aavegotchi.AddressAavegotchi, event.Erc721TokenAddress, event.Erc721TokenId, nil)
}

func (w *worker) handleERC721ExecutedListing(ctx context.Context, task *source.Task, log ethereum.Log, feed *schema.Feed) (*schema.Action, error) {
	event, err := w.erc721MarketplaceFilterer.ParseERC721ExecutedListing(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse erc721 executed listing: %w", err)
	}

	feed.Type = filter.TypeMetaverseTrade

	metadataType := lo.If(feed.From == event.Buyer.String(), metadata.ActionMetaverseTradeBuy).Else(metadata.ActionMetaverseTradeSell)

	return w.buildTradeAction(ctx, log.BlockNumber, task.ChainID, feed.Type, metadataType, event.Seller, event.Buyer, event.Erc721TokenAddress, event.Erc721TokenId, nil)
}

func (w *worker) handleERC1155ListingAdd(ctx context.Context, task *source.Task, log ethereum.Log, feed *schema.Feed) (*schema.Action, error) {
	event, err := w.erc1155MarketplaceFilterer.ParseERC1155ListingAdd(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse erc1155 listing add: %w", err)
	}

	feed.Type = filter.TypeMetaverseTrade

	return w.buildTradeAction(ctx, log.BlockNumber, task.ChainID, feed.Type, metadata.ActionMetaverseTradeList, event.Seller, aavegotchi.AddressAavegotchi, event.Erc1155TokenAddress, event.Erc1155TypeId, big.NewInt(1))
}

func (w *worker) handleERC1155ExecutedListing(ctx context.Context, task *source.Task, log ethereum.Log, feed *schema.Feed) (*schema.Action, error) {
	event, err := w.erc1155MarketplaceFilterer.ParseERC1155ExecutedListing(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse erc1155 executed listing: %w", err)
	}

	feed.Type = filter.TypeMetaverseTrade

	metadataType := lo.If(feed.From == event.Buyer.String(), metadata.ActionMetaverseTradeBuy).Else(metadata.ActionMetaverseTradeSell)

	return w.buildTradeAction(ctx, log.BlockNumber, task.ChainID, feed.Type, metadataType, event.Seller, event.Buyer, event.Erc1155TokenAddress, event.Erc1155TypeId, big.NewInt(1))
}

func (w *worker) handleERC20TransferLog(ctx context.Context, task *source.Task, log ethereum.Log, feed *schema.Feed) (*schema.Action, error) {
	event, err := w.erc20Filterer.ParseTransfer(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse erc20 transfer: %w", err)
	}

	actionType := lo.If(event.To == ethereum.AddressGenesis, filter.TypeMetaverseBurn).
		ElseIf(event.From == ethereum.AddressGenesis, filter.TypeMetaverseMint).
		Else(filter.TypeMetaverseTransfer)

	if feed.Type == filter.TypeUnknown {
		feed.Type = actionType
	}

	return w.buildTransferAction(ctx, log.BlockNumber, task.ChainID, actionType, event.From, event.To, event.Raw.Address, nil, event.Value)
}

func (w *worker) handleMetaverseTradeCost(_ context.Context, feed *schema.Feed) (*schema.Feed, error) {
	// count the cost of the trade
	cost := metadata.Token{
		Value: lo.ToPtr(decimal.NewFromInt(0)),
	}

	for _, action := range feed.Actions {
		if action.From != feed.From {
			continue
		}

		if action.Type == filter.TypeMetaverseTransfer && action.Metadata.(metadata.MetaverseTransfer).Value != nil {
			if cost.Value.IsZero() {
				cost = metadata.Token(action.Metadata.(metadata.MetaverseTransfer))
			} else {
				cost.Value = lo.ToPtr(cost.Value.Add(decimal.NewFromBigInt(action.Metadata.(metadata.MetaverseTransfer).Value.BigInt(), 0)))
			}
		}
	}

	for _, action := range feed.Actions {
		if action.Type == feed.Type {
			action.Metadata = metadata.MetaverseTrade{
				Action: action.Metadata.(metadata.MetaverseTrade).Action,
				Token:  action.Metadata.(metadata.MetaverseTrade).Token,
				Cost:   cost,
			}

			feed.Actions = []*schema.Action{action}

			return feed, nil
		}
	}

	return nil, fmt.Errorf("no metaverse trade action found")
}

func (w *worker) buildTradeAction(
	ctx context.Context,
	blockNumber *big.Int,
	chainID uint64,
	actionType filter.Type,
	metadataAction metadata.MetaverseTradeAction,
	from, to, tokenAddress common.Address,
	tokenID, value *big.Int,
) (*schema.Action, error) {
	tokenMetadata, err := w.tokenClient.Lookup(ctx, chainID, lo.ToPtr(tokenAddress), tokenID, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup token: %w", err)
	}

	tokenMetadata.Value = lo.If[*decimal.Decimal](value == nil, nil).ElseF(func() *decimal.Decimal {
		return lo.ToPtr(decimal.NewFromBigInt(value, 0))
	})

	return &schema.Action{
		Type:     actionType,
		Platform: filter.PlatformAavegotchi.String(),
		From:     from.String(),
		To:       to.String(),
		Metadata: metadata.MetaverseTrade{
			Action: metadataAction,
			Token:  *tokenMetadata,
		},
	}, nil
}

func (w *worker) buildTransferAction(
	ctx context.Context,
	blockNumber *big.Int,
	chainID uint64,
	actionType filter.Type,
	from, to, tokenAddress common.Address,
	tokenID, value *big.Int,
) (*schema.Action, error) {
	tokenMetadata, err := w.tokenClient.Lookup(ctx, chainID, lo.ToPtr(tokenAddress), tokenID, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup token: %w", err)
	}

	tokenMetadata.Value = lo.If[*decimal.Decimal](value == nil, nil).ElseF(func() *decimal.Decimal {
		return lo.ToPtr(decimal.NewFromBigInt(value, 0))
	})

	return &schema.Action{
		Type:     actionType,
		Platform: filter.PlatformAavegotchi.String(),
		From:     from.String(),
		To:       lo.If(to == ethereum.AddressGenesis, "").Else(to.String()),
		Metadata: metadata.MetaverseTransfer(*tokenMetadata),
	}, nil
}

func NewWorker(config *engine.Config) (engine.Worker, error) {
	ethereumClient, err := ethereum.Dial(context.Background(), config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("dial Ethereum: %w", err)
	}

	return &worker{
		tokenClient:                token.NewClient(ethereumClient),
		erc20Filterer:              lo.Must(erc20.NewERC20Filterer(ethereum.AddressGenesis, nil)),
		erc721MarketplaceFilterer:  lo.Must(aavegotchi.NewERC721MarketplaceFilterer(aavegotchi.AddressAavegotchi, nil)),
		erc1155MarketplaceFilterer: lo.Must(aavegotchi.NewERC1155MarketplaceFilterer(aavegotchi.AddressAavegotchi, nil)),
	}, nil
}
