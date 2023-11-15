// Code generated by "enumer --values --type=ChainEthereum --linecomment --output chain_ethereum_string.go --json --sql"; DO NOT EDIT.

package filter

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _ChainEthereumName = "mainnet"

var _ChainEthereumIndex = [...]uint8{0, 7}

const _ChainEthereumLowerName = "mainnet"

func (i ChainEthereum) String() string {
	i -= 1
	if i >= ChainEthereum(len(_ChainEthereumIndex)-1) {
		return fmt.Sprintf("ChainEthereum(%d)", i+1)
	}
	return _ChainEthereumName[_ChainEthereumIndex[i]:_ChainEthereumIndex[i+1]]
}

func (ChainEthereum) Values() []string {
	return ChainEthereumStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ChainEthereumNoOp() {
	var x [1]struct{}
	_ = x[ChainEthereumMainnet-(1)]
}

var _ChainEthereumValues = []ChainEthereum{ChainEthereumMainnet}

var _ChainEthereumNameToValueMap = map[string]ChainEthereum{
	_ChainEthereumName[0:7]:      ChainEthereumMainnet,
	_ChainEthereumLowerName[0:7]: ChainEthereumMainnet,
}

var _ChainEthereumNames = []string{
	_ChainEthereumName[0:7],
}

// ChainEthereumString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ChainEthereumString(s string) (ChainEthereum, error) {
	if val, ok := _ChainEthereumNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ChainEthereumNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ChainEthereum values", s)
}

// ChainEthereumValues returns all values of the enum
func ChainEthereumValues() []ChainEthereum {
	return _ChainEthereumValues
}

// ChainEthereumStrings returns a slice of all String values of the enum
func ChainEthereumStrings() []string {
	strs := make([]string, len(_ChainEthereumNames))
	copy(strs, _ChainEthereumNames)
	return strs
}

// IsAChainEthereum returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ChainEthereum) IsAChainEthereum() bool {
	for _, v := range _ChainEthereumValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ChainEthereum
func (i ChainEthereum) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ChainEthereum
func (i *ChainEthereum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ChainEthereum should be a string, got %s", data)
	}

	var err error
	*i, err = ChainEthereumString(s)
	return err
}

func (i ChainEthereum) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *ChainEthereum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of ChainEthereum: %[1]T(%[1]v)", value)
	}

	val, err := ChainEthereumString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
