package filter

var _ Chain = (*ChainEthereum)(nil)

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=ChainArweave --linecomment --output chain_arweave_string.go --json --sql
type ChainArweave uint64

//goland:noinspection GoMixedReceiverTypes
func (i ChainArweave) Network() Network {
	return NetworkEthereum
}

//goland:noinspection GoMixedReceiverTypes
func (i ChainArweave) ID() uint64 {
	return uint64(i)
}

//goland:noinspection GoMixedReceiverTypes
func (i ChainArweave) Name() string {
	return i.String()
}

const (
	ChainArweaveMainnet ChainArweave = 1 // mainnet
)
