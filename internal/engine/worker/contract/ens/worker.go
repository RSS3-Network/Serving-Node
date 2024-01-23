package ens

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/naturalselectionlabs/rss3-node/config"
	"github.com/naturalselectionlabs/rss3-node/internal/database"
	"github.com/naturalselectionlabs/rss3-node/internal/database/model"
	"github.com/naturalselectionlabs/rss3-node/internal/engine"
	source "github.com/naturalselectionlabs/rss3-node/internal/engine/source/ethereum"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/ens"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc1155"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc20"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/contract/erc721"
	"github.com/naturalselectionlabs/rss3-node/provider/ethereum/token"
	"github.com/naturalselectionlabs/rss3-node/schema"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/naturalselectionlabs/rss3-node/schema/metadata"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// Worker is the worker for ENS.
var _ engine.Worker = (*worker)(nil)

type worker struct {
	// Base
	config         *config.Module
	ethereumClient ethereum.Client
	tokenClient    token.Client
	databaseClient database.Client
	// Specific filters
	baseRegistrarImplementationFilterer *ens.BaseRegistrarImplementationFilterer
	ethRegistrarControllerV1Filterer    *ens.ETHRegistrarControllerV1Filterer
	ethRegistrarControllerV2Filterer    *ens.ETHRegistrarControllerV2Filterer
	publicResolverV1Filterer            *ens.PublicResolverV1Filterer
	publicResolverV2Filterer            *ens.PublicResolverV2Filterer
	nameWrapperFilterer                 *ens.NameWrapperFilterer
	// Common filters
	erc20Filterer   *erc20.ERC20Filterer
	erc721Filterer  *erc721.ERC721Filterer
	erc1155Filterer *erc1155.ERC1155Filterer
}

func (w *worker) Name() string {
	return filter.ENS.String()
}

// Filter ens contract address and event hash.
func (w *worker) Filter() engine.SourceFilter {
	return &source.Filter{
		LogAddresses: []common.Address{
			ens.AddressBaseRegistrarImplementation,
			ens.AddressETHRegistrarControllerV1,
			ens.AddressETHRegistrarControllerV2,
			ens.AddressPublicResolverV1,
			ens.AddressPublicResolverV2,
			ens.AddressNameWrapper,
		},
		LogTopics: []common.Hash{
			ens.EventNameRegisteredControllerV1,
			ens.EventNameRegisteredControllerV2,
			ens.EventNameRenewed,
			ens.EventTextChanged,
			ens.EventTextChangedWithValue,
			ens.EventNameWrapped,
			ens.EventNameUnwrapped,

			ens.EventFusesSet,
			ens.EventAddressChanged,
			ens.EventContenthashChanged,
			ens.EventNameChanged,
			ens.EventPubkeyChanged,
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

	// Build default ens feed from task.
	feed, err := ethereumTask.BuildFeed(schema.WithFeedPlatform(filter.PlatformENS))
	if err != nil {
		return nil, fmt.Errorf("build feed: %w", err)
	}

	exist := lo.ContainsBy(ethereumTask.Receipt.Logs, func(log *ethereum.Log) bool {
		return w.matchEnsNameRegisteredV2(ctx, *log) || w.matchEnsNameRegisteredV1(ctx, *log)
	})

	// Match and handle ethereum logs.
	for _, log := range ethereumTask.Receipt.Logs {
		var (
			actions []*schema.Action
			err     error
		)

		if len(log.Topics) == 0 {
			continue
		}

		if exist {
			switch {
			case w.matchEnsNameRegisteredV1(ctx, *log):
				feed.Type = filter.TypeCollectibleTrade
				actions, err = w.transformEnsNameRegisteredV1(ctx, *log, ethereumTask)
			case w.matchEnsNameRegisteredV2(ctx, *log):
				feed.Type = filter.TypeCollectibleTrade
				actions, err = w.transformEnsNameRegisteredV2(ctx, *log, ethereumTask)
			default:
				continue
			}
		} else {
			switch {
			case w.matchEnsNameRenewed(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsNameRenewed(ctx, *log, ethereumTask)
			case w.matchEnsTextChanged(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsTextChanged(ctx, *log, ethereumTask)
			case w.matchEnsTextChangedWithValue(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsTextChangedWithValue(ctx, *log, ethereumTask)
			case w.matchEnsNameWrapped(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsNameWrapped(ctx, *log, ethereumTask)
			case w.matchEnsNameUnwrapped(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsNameUnwrapped(ctx, *log, ethereumTask)
			case w.matchEnsFusesSet(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsFusesSet(ctx, *log, ethereumTask)
			case w.matchEnsContenthashChanged(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsContenthashChanged(ctx, *log, ethereumTask)
			case w.matchEnsNameChanged(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsNameChanged(ctx, *log, ethereumTask)
			case w.matchEnsAddressChanged(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsAddressChanged(ctx, *log, ethereumTask)
			case w.matchEnsPubkeyChanged(ctx, *log):
				feed.Type = filter.TypeSocialProfile
				actions, err = w.transformEnsPubkeyChanged(ctx, *log, ethereumTask)
			default:
				continue
			}
		}

		if err != nil {
			return nil, err
		}

		// Change feed type to the first action type.
		for _, action := range actions {
			feed.Type = action.Type
		}

		feed.Actions = append(feed.Actions, actions...)
	}

	return feed, nil
}

// NewWorker creates a new ENS worker.
func NewWorker(config *config.Module, databaseClient database.Client) (engine.Worker, error) {
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

	// Initialize database client.
	instance.databaseClient = databaseClient

	// Initialize ens filterers and callers.
	instance.baseRegistrarImplementationFilterer = lo.Must(ens.NewBaseRegistrarImplementationFilterer(ethereum.AddressGenesis, nil))
	instance.ethRegistrarControllerV1Filterer = lo.Must(ens.NewETHRegistrarControllerV1Filterer(ethereum.AddressGenesis, nil))
	instance.ethRegistrarControllerV2Filterer = lo.Must(ens.NewETHRegistrarControllerV2Filterer(ethereum.AddressGenesis, nil))
	instance.publicResolverV1Filterer = lo.Must(ens.NewPublicResolverV1Filterer(ethereum.AddressGenesis, nil))
	instance.publicResolverV2Filterer = lo.Must(ens.NewPublicResolverV2Filterer(ethereum.AddressGenesis, nil))
	instance.nameWrapperFilterer = lo.Must(ens.NewNameWrapperFilterer(ethereum.AddressGenesis, nil))

	instance.erc20Filterer = lo.Must(erc20.NewERC20Filterer(ethereum.AddressGenesis, nil))
	instance.erc721Filterer = lo.Must(erc721.NewERC721Filterer(ethereum.AddressGenesis, nil))
	instance.erc1155Filterer = lo.Must(erc1155.NewERC1155Filterer(ethereum.AddressGenesis, nil))

	return &instance, nil
}

func (w *worker) matchEnsNameRegisteredV1(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressETHRegistrarControllerV1 && log.Topics[0] == ens.EventNameRegisteredControllerV1
}

func (w *worker) matchEnsNameRegisteredV2(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressETHRegistrarControllerV2 && log.Topics[0] == ens.EventNameRegisteredControllerV2
}

func (w *worker) matchEnsNameRenewed(_ context.Context, log ethereum.Log) bool {
	return (log.Address == ens.AddressETHRegistrarControllerV1 || log.Address == ens.AddressETHRegistrarControllerV2) && log.Topics[0] == ens.EventNameRenewed
}

func (w *worker) matchEnsTextChanged(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressPublicResolverV1 && log.Topics[0] == ens.EventTextChanged
}

func (w *worker) matchEnsTextChangedWithValue(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressPublicResolverV2 && log.Topics[0] == ens.EventTextChangedWithValue
}

func (w *worker) matchEnsNameWrapped(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressNameWrapper && log.Topics[0] == ens.EventNameWrapped
}

func (w *worker) matchEnsNameUnwrapped(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressNameWrapper && log.Topics[0] == ens.EventNameUnwrapped
}

func (w *worker) matchEnsFusesSet(_ context.Context, log ethereum.Log) bool {
	return log.Address == ens.AddressNameWrapper && log.Topics[0] == ens.EventFusesSet
}

func (w *worker) matchEnsContenthashChanged(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == ens.EventContenthashChanged
}

func (w *worker) matchEnsNameChanged(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == ens.EventNameChanged
}

func (w *worker) matchEnsAddressChanged(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == ens.EventAddressChanged
}

func (w *worker) matchEnsPubkeyChanged(_ context.Context, log ethereum.Log) bool {
	return log.Topics[0] == ens.EventPubkeyChanged
}

func (w *worker) transformEnsNameRegisteredV1(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.ethRegistrarControllerV1Filterer.ParseNameRegistered(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseNameRegistered event: %w", err)
	}

	action, err := w.buildEthereumENSRegisterAction(ctx, task, event.Label, ethereum.AddressGenesis, event.Owner, event.Cost, event.Name)

	if err != nil {
		return nil, fmt.Errorf("build collectible trade action: %w", err)
	}

	return []*schema.Action{action}, nil
}

func (w *worker) transformEnsNameRegisteredV2(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.ethRegistrarControllerV2Filterer.ParseNameRegistered(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseNameRegistered event: %w", err)
	}

	action, err := w.buildEthereumENSRegisterAction(ctx, task, event.Label, ethereum.AddressGenesis, event.Owner, event.BaseCost, event.Name)

	if err != nil {
		return nil, fmt.Errorf("build collectible trade action: %w", err)
	}

	return []*schema.Action{action}, nil
}

func (w *worker) transformEnsNameRenewed(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.ethRegistrarControllerV2Filterer.ParseNameRenewed(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseNameRenewed event: %w", err)
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, event.Expires,
			fmt.Sprintf("%s.eth", event.Name), "", "",
			metadata.ActionSocialProfileRenew),
	}, nil
}

func (w *worker) transformEnsTextChanged(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.publicResolverV1Filterer.ParseTextChanged(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseTextChanged event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, event.Key, "",
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) transformEnsTextChangedWithValue(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.publicResolverV2Filterer.ParseTextChanged(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseTextChangedWithValue event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, event.Key, event.Value,
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) transformEnsNameWrapped(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.nameWrapperFilterer.ParseNameWrapped(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseNameWrapped event: %w", err)
	}

	var name string

	decodeEnsName, err := w.decodeName(event.Name)
	if err != nil {
		return nil, err
	}

	if len(decodeEnsName) != 2 {
		return nil, nil
	}

	name = decodeEnsName[1]

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, "", "",
			metadata.ActionSocialProfileWrap),
	}, nil
}

func (w *worker) transformEnsNameUnwrapped(ctx context.Context, log ethereum.Log, _ *source.Task) ([]*schema.Action, error) {
	event, err := w.nameWrapperFilterer.ParseNameUnwrapped(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseNameUnwrapped event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, event.Owner, log.Address, nil,
			name, "", "",
			metadata.ActionSocialProfileUnwrap),
	}, nil
}

func (w *worker) transformEnsFusesSet(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.nameWrapperFilterer.ParseFusesSet(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseFusesSet event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, ens.FusesSet.String(), strconv.Itoa(int(event.Fuses)),
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) transformEnsContenthashChanged(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.publicResolverV2Filterer.ParseContenthashChanged(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseContenthashChanged event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, ens.ContentHashChanged.String(), common.BytesToHash(event.Hash).Hex(),
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) transformEnsNameChanged(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.publicResolverV2Filterer.ParseNameChanged(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseNameChanged event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, ens.NameChanged.String(), event.Name,
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) transformEnsAddressChanged(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.publicResolverV2Filterer.ParseAddressChanged(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParseAddressChanged event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, ens.AddressChanged.String(), common.BytesToAddress(event.NewAddress).Hex(),
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) transformEnsPubkeyChanged(ctx context.Context, log ethereum.Log, task *source.Task) ([]*schema.Action, error) {
	event, err := w.publicResolverV2Filterer.ParsePubkeyChanged(log.Export())
	if err != nil {
		return nil, fmt.Errorf("ParsePubkeyChanged event: %w", err)
	}

	name, err := w.getEnsName(ctx, log.BlockNumber, event.Node)
	if err != nil {
		return nil, err
	}

	return []*schema.Action{
		w.buildEthereumENSProfileAction(ctx, task.Transaction.From, log.Address, nil,
			name, ens.PubkeyChanged.String(), common.Bytes2Hex(secp256k1.CompressPubkey(new(big.Int).SetBytes(event.X[:]), new(big.Int).SetBytes(event.Y[:]))),
			metadata.ActionSocialProfileUpdate),
	}, nil
}

func (w *worker) buildEthereumENSRegisterAction(ctx context.Context, task *source.Task, labels [32]byte, from, to common.Address, cost *big.Int, name string) (*schema.Action, error) {
	tokenMetadata, err := w.tokenClient.Lookup(ctx, task.ChainID, lo.ToPtr(ens.AddressBaseRegistrarImplementation), new(big.Int).SetBytes(labels[:]), task.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("lookup token metadata %s %s: %w", ens.AddressBaseRegistrarImplementation.String(), new(big.Int).SetBytes(labels[:]), err)
	}

	tokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(big.NewInt(1), 0))

	costTokenMetadata, err := w.tokenClient.Lookup(ctx, task.ChainID, nil, nil, task.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("lookup token metadata %s: %w", "", err)
	}

	costTokenMetadata.Value = lo.ToPtr(decimal.NewFromBigInt(cost, 0))

	// Save namehash into database for further query requirements
	fullName := fmt.Sprintf("%s.%s", name, "eth")
	if err = w.databaseClient.SaveDatasetENSNamehash(ctx, &model.ENSNamehash{
		Name: fullName,
		Hash: ens.NameHash(fullName),
	}); err != nil {
		return nil, fmt.Errorf("save dataset ens namehash: %w", err)
	}

	return &schema.Action{
		Type:     filter.TypeCollectibleTrade,
		Platform: filter.PlatformENS.String(),
		From:     from.String(),
		To:       to.String(),
		Metadata: metadata.CollectibleTrade{
			Action: metadata.ActionCollectibleTradeBuy,
			Token:  *tokenMetadata,
			Cost:   costTokenMetadata,
		},
	}, nil
}

func (w *worker) buildEthereumENSProfileAction(_ context.Context, from, to common.Address, expires *big.Int, name, key, value string, action metadata.SocialProfileAction) *schema.Action {
	label := strings.Split(name, ".eth")[0]
	tokenID := common.BytesToHash(crypto.Keccak256([]byte(label))).Big()

	socialProfile := metadata.SocialProfile{
		Action:    action,
		ProfileID: tokenID.String(),
		Address:   ens.AddressBaseRegistrarImplementation,
		Handle:    name,
		ImageURI:  fmt.Sprintf("https://metadata.ens.domains/mainnet/%s/%s/image", ens.AddressBaseRegistrarImplementation.String(), tokenID.String()),
		Key:       key,
		Value:     value,
	}

	if expires != nil {
		socialProfile.Expiry = lo.ToPtr(time.Unix(expires.Int64(), 0))
	}

	return &schema.Action{
		Type:     filter.TypeSocialProfile,
		Platform: filter.PlatformENS.String(),
		From:     from.String(),
		To:       to.String(),
		Metadata: socialProfile,
	}
}