package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/internal/database/model"
)

var _ model.ENSNamehashTransformer = (*DatasetENSNamehash)(nil)

type DatasetENSNamehash struct {
	HashHex string `gorm:"column:hash_hex;primary_key"`
	Name    string `gorm:"column:name;index:idx_ensnamehash_name"`
}

func (d *DatasetENSNamehash) Import(ensNamehash *model.ENSNamehash) (err error) {
	d.HashHex = ensNamehash.Hash.Hex()
	d.Name = ensNamehash.Name

	return nil
}

func (d *DatasetENSNamehash) Export() (*model.ENSNamehash, error) {
	profile := model.ENSNamehash{
		Hash: common.HexToHash(d.HashHex),
		Name: d.Name,
	}

	return &profile, nil
}
