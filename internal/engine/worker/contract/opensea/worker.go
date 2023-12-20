package opensea

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/internal/engine"
	source "github.com/naturalselectionlabs/rss3-node/internal/engine/source/ethereum"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc1155"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc20"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc721"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/opensea"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/token"
	"github.com/naturalselectionlabs/rss3-node/schema"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/naturalselectionlabs/rss3-node/schema/metadata"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var _ engine.Worker = (*worker)(nil)

type worker struct {
	config                   *engine.Config
	ethereumClient           ethereum.Client
	tokenClient              token.Client
	wyvernExchangeV1Filterer *opensea.WyvernExchangeV1Filterer
	wyvernExchangeV2Filterer *opensea.WyvernExchangeV2Filterer
	seaportV1Filterer        *opensea.SeaportV1Filterer
	erc20Filterer            *erc20.ERC20Filterer
	erc721Filterer           *erc721.ERC721Filterer
	erc1155Filterer          *erc1155.ERC1155Filterer
}

func (w *worker) Name() string {
	return engine.OpenSea.String()
}

func (w *worker) Filter() engine.SourceFilter {
	return &source.Filter{
		LogAddresses: []common.Address{
			opensea.AddressWyvernExchangeV1,
			opensea.AddressWyvernExchangeV2,
			opensea.AddressSeaportV1Dot1,
			opensea.AddressSeaportV1Dot2,
			opensea.AddressSeaportV1Dot3,
			opensea.AddressSeaportV1Dot4,
			opensea.AddressSeaportV1Dot5,
		},
		LogTopics: []common.Hash{
			opensea.EventHashWyvernExchangeV1OrdersMatched,
			opensea.EventHashWyvernExchangeV2OrdersMatched,
			opensea.EventHashSeaportV1OrderFulfilled,
		},
	}
}

func (w *worker) Match(_ context.Context, task engine.Task) (bool, error) {
	switch task.GetNetwork() {
	// Match Ethereum network.
	case filter.NetworkEthereum:
		// Match Ethereum task.
		task := task.(*source.Task)
		if task.Transaction.To == nil {
			return false, nil
		}

		// Match OpenSea contract.
		switch *task.Transaction.To {
		case
			opensea.AddressWyvernExchangeV1,
			opensea.AddressWyvernExchangeV2,
			opensea.AddressSeaportV1Dot1,
			opensea.AddressSeaportV1Dot2,
			opensea.AddressSeaportV1Dot3,
			opensea.AddressSeaportV1Dot4,
			opensea.AddressSeaportV1Dot5:
			return true, nil
		default:
			return false, nil
		}
	default:
		return false, nil
	}
}

func (w *worker) Transform(ctx context.Context, task engine.Task) (*schema.Feed, error) {
	ethereumTask, ok := task.(*source.Task)
	if !ok {
		return nil, fmt.Errorf("invalid task type: %T", task)
	}

	feed, err := ethereumTask.BuildFeed(schema.WithFeedPlatform(filter.PlatformOpenSea))
	if err != nil {
		return nil, fmt.Errorf("build feed: %w", err)
	}

	// Match and handle logs.
	for _, log := range ethereumTask.Receipt.Logs {
		var (
			actions []*schema.Action
			err     error
		)

		switch {
		case w.matchWyvernExchangeV1Orders(ethereumTask, log):
			actions, err = w.transformWyvernExchangeV1Orders(ctx, ethereumTask, log)
		case w.matchWyvernExchangeV2Orders(ethereumTask, log):
			actions, err = w.transformWyvernExchangeV2Orders(ctx, ethereumTask, log)
		case w.matchSeaportV1OrderFulfilled(ethereumTask, log):
			actions, err = w.transformSeaportV1OrderFulfilled(ctx, ethereumTask, log)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		// Overwrite the type for feed.
		for _, action := range actions {
			feed.Type = action.Type
		}

		feed.Actions = append(feed.Actions, actions...)
	}

	return feed, nil
}

func (w *worker) matchWyvernExchangeV1Orders(_ *source.Task, log *ethereum.Log) bool {
	return log.Address == opensea.AddressWyvernExchangeV1 && len(log.Topics) == 4 && contract.MatchEventHashes(log.Topics[0], opensea.EventHashWyvernExchangeV1OrdersMatched)
}

func (w *worker) matchWyvernExchangeV2Orders(_ *source.Task, log *ethereum.Log) bool {
	return log.Address == opensea.AddressWyvernExchangeV2 && len(log.Topics) == 4 && contract.MatchEventHashes(log.Topics[0], opensea.EventHashWyvernExchangeV2OrdersMatched)
}

func (w *worker) matchSeaportV1OrderFulfilled(_ *source.Task, log *ethereum.Log) bool {
	if !contract.MatchEventHashes(log.Topics[0], opensea.EventHashSeaportV1OrderFulfilled) {
		return false
	}

	return contract.MatchAddresses(
		log.Address,
		opensea.AddressSeaportV1Dot0, opensea.AddressSeaportV1Dot1, opensea.AddressSeaportV1Dot2,
		opensea.AddressSeaportV1Dot3, opensea.AddressSeaportV1Dot4, opensea.AddressSeaportV1Dot5,
	)
}

func (w *worker) transformWyvernExchangeV1Orders(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	event, err := w.wyvernExchangeV1Filterer.ParseOrdersMatched(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse OrdersMatched event: %w", err)
	}

	var (
		nftAddress common.Address
		nftID      *big.Int
		nftValue   *big.Int
	)

	for _, log := range task.Receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		switch {
		case log.Topics[0] == erc721.EventHashTransfer && len(log.Topics) == 4:
			transferEvent, err := w.erc721Filterer.ParseTransfer(log.Export())
			if err != nil {
				return nil, fmt.Errorf("parse transfer event: %w", err)
			}

			if event.Taker != transferEvent.To {
				event.Maker, event.Taker = event.Taker, event.Maker
			}

			nftAddress = transferEvent.Raw.Address
			nftID = transferEvent.TokenId
			nftValue = big.NewInt(1)
		case log.Topics[0] == erc1155.EventHashTransferSingle:
			transferEvent, err := w.erc1155Filterer.ParseTransferSingle(log.Export())
			if err != nil {
				return nil, fmt.Errorf("parse transfer event: %w", err)
			}

			if event.Taker != transferEvent.To {
				event.Maker, event.Taker = event.Taker, event.Maker
			}

			nftAddress = transferEvent.Raw.Address
			nftID = transferEvent.Id
			nftValue = transferEvent.Value
		default:
			continue
		}
	}

	action, err := w.buildEthereumCollectibleTradeAction(ctx, task, event.Maker, event.Taker, nftAddress, nftID, nftValue, nil, event.Price)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		action,
	}, nil
}

func (w *worker) transformWyvernExchangeV2Orders(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	return w.transformWyvernExchangeV1Orders(ctx, task, log)
}

func (w *worker) transformSeaportV1OrderFulfilled(ctx context.Context, task *source.Task, log *ethereum.Log) ([]*schema.Action, error) {
	event, err := w.seaportV1Filterer.ParseOrderFulfilled(log.Export())
	if err != nil {
		return nil, fmt.Errorf("parse OrderFulfilled event: %w", err)
	}

	seller, buyer := event.Offerer, event.Recipient

	var (
		offerToken      *common.Address
		offerValue      *big.Int
		nft             common.Address
		nftID, nftValue *big.Int
	)

	for _, offer := range event.Offer {
		switch offer.ItemType {
		case opensea.ItemTypeNative:
			offerToken, offerValue = nil, offer.Amount
		case opensea.ItemTypeERC20:
			offerToken, offerValue = lo.ToPtr(offer.Token), offer.Amount
		case opensea.ItemTypeERC721, opensea.ItemTypeERC1155:
			nft, nftID, nftValue = offer.Token, offer.Identifier, offer.Amount
		}
	}

	for _, consideration := range event.Consideration {
		switch consideration.ItemType {
		case opensea.ItemTypeNative:
			// Ignore token transfers where receipt is not seller.
			if consideration.Recipient != seller {
				continue
			}

			offerToken, offerValue = nil, consideration.Amount
		case opensea.ItemTypeERC20:
			// Ignore token transfers where receipt is not seller.
			if consideration.Recipient != seller {
				continue
			}

			offerToken, offerValue = lo.ToPtr(consideration.Token), consideration.Amount
		case opensea.ItemTypeERC721, opensea.ItemTypeERC1155:
			// The offerer is the seller, and the recipient is the buyer.
			seller, buyer = buyer, seller

			nft, nftID, nftValue = consideration.Token, consideration.Identifier, consideration.Amount
		}
	}

	action, err := w.buildEthereumCollectibleTradeAction(ctx, task, seller, buyer, nft, nftID, nftValue, offerToken, offerValue)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		action,
	}, nil
}

func (w *worker) buildEthereumCollectibleTradeAction(ctx context.Context, task *source.Task, seller, buyer, nft common.Address, nftID, nftValue *big.Int, offerToken *common.Address, offerValue *big.Int) (*schema.Action, error) {
	if nftID == nil {
		return nil, fmt.Errorf("nft id is nil")
	}

	nftValue = lo.Ternary(nftValue == nil, big.NewInt(1), nftValue)

	collectibleTokenMetadata, err := w.tokenClient.Lookup(ctx, task.ChainID, &nft, nftID, task.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("lookup collectible token metadata: %w", err)
	}

	collectibleTokenMetadata.ID = lo.ToPtr(decimal.NewFromBigInt(nftID, 0))
	collectibleTokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(nftValue, 0))

	costTokenMetadata, err := w.tokenClient.Lookup(ctx, task.ChainID, offerToken, nil, task.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("lookup collectible token metadata: %w", err)
	}

	if offerValue == nil {
		offerValue = big.NewInt(0)
	}

	costTokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(offerValue, 0))

	action := schema.Action{
		Type:     filter.TypeCollectibleTrade,
		Platform: filter.PlatformOpenSea.String(),
		From:     seller.String(),
		To:       buyer.String(),
		Metadata: metadata.CollectibleTrade{
			Action: lo.If(task.Transaction.From == seller,
				metadata.ActionCollectibleTradeSell).Else(metadata.ActionCollectibleTradeBuy),
			Token: *collectibleTokenMetadata,
			Cost:  costTokenMetadata,
		},
	}

	return &action, nil
}

// NewWorker creates a new RSS3 worker.
func NewWorker(config *engine.Config) (engine.Worker, error) {
	var (
		err      error
		instance = worker{
			config: config,
		}
	)

	if instance.ethereumClient, err = ethereum.Dial(context.Background(), config.Endpoint); err != nil {
		return nil, fmt.Errorf("initialize ethereum client: %w", err)
	}

	instance.tokenClient = token.NewClient(instance.ethereumClient)

	// Initialize opensea filterers.
	instance.wyvernExchangeV1Filterer = lo.Must(opensea.NewWyvernExchangeV1Filterer(ethereum.AddressGenesis, nil))
	instance.wyvernExchangeV2Filterer = lo.Must(opensea.NewWyvernExchangeV2Filterer(ethereum.AddressGenesis, nil))
	instance.seaportV1Filterer = lo.Must(opensea.NewSeaportV1Filterer(ethereum.AddressGenesis, nil))
	instance.erc20Filterer = lo.Must(erc20.NewERC20Filterer(ethereum.AddressGenesis, nil))
	instance.erc721Filterer = lo.Must(erc721.NewERC721Filterer(ethereum.AddressGenesis, nil))
	instance.erc1155Filterer = lo.Must(erc1155.NewERC1155Filterer(ethereum.AddressGenesis, nil))

	return &instance, nil
}
