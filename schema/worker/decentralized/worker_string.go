// Code generated by "enumer --values --type=Worker --linecomment --output worker_string.go --json --yaml --sql"; DO NOT EDIT.

package decentralized

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _WorkerName = "aaveaavegotchiarbitrumcorecrossbellcurveenshighlightiqwikikiwistandlenslidolooksraremattersmirrormomoka1inchopenseaoptimismparagraphrss3savmstargateuniswapvsl"

var _WorkerIndex = [...]uint8{0, 4, 14, 22, 26, 35, 40, 43, 52, 58, 67, 71, 75, 84, 91, 97, 103, 108, 115, 123, 132, 136, 140, 148, 155, 158}

const _WorkerLowerName = "aaveaavegotchiarbitrumcorecrossbellcurveenshighlightiqwikikiwistandlenslidolooksraremattersmirrormomoka1inchopenseaoptimismparagraphrss3savmstargateuniswapvsl"

func (i Worker) String() string {
	i -= 1
	if i < 0 || i >= Worker(len(_WorkerIndex)-1) {
		return fmt.Sprintf("Worker(%d)", i+1)
	}
	return _WorkerName[_WorkerIndex[i]:_WorkerIndex[i+1]]
}

func (Worker) Values() []string {
	return WorkerStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _WorkerNoOp() {
	var x [1]struct{}
	_ = x[Aave-(1)]
	_ = x[Aavegotchi-(2)]
	_ = x[Arbitrum-(3)]
	_ = x[Core-(4)]
	_ = x[Crossbell-(5)]
	_ = x[Curve-(6)]
	_ = x[ENS-(7)]
	_ = x[Highlight-(8)]
	_ = x[IQWiki-(9)]
	_ = x[KiwiStand-(10)]
	_ = x[Lens-(11)]
	_ = x[Lido-(12)]
	_ = x[Looksrare-(13)]
	_ = x[Matters-(14)]
	_ = x[Mirror-(15)]
	_ = x[Momoka-(16)]
	_ = x[Oneinch-(17)]
	_ = x[OpenSea-(18)]
	_ = x[Optimism-(19)]
	_ = x[Paragraph-(20)]
	_ = x[RSS3-(21)]
	_ = x[SAVM-(22)]
	_ = x[Stargate-(23)]
	_ = x[Uniswap-(24)]
	_ = x[VSL-(25)]
}

var _WorkerValues = []Worker{Aave, Aavegotchi, Arbitrum, Core, Crossbell, Curve, ENS, Highlight, IQWiki, KiwiStand, Lens, Lido, Looksrare, Matters, Mirror, Momoka, Oneinch, OpenSea, Optimism, Paragraph, RSS3, SAVM, Stargate, Uniswap, VSL}

var _WorkerNameToValueMap = map[string]Worker{
	_WorkerName[0:4]:          Aave,
	_WorkerLowerName[0:4]:     Aave,
	_WorkerName[4:14]:         Aavegotchi,
	_WorkerLowerName[4:14]:    Aavegotchi,
	_WorkerName[14:22]:        Arbitrum,
	_WorkerLowerName[14:22]:   Arbitrum,
	_WorkerName[22:26]:        Core,
	_WorkerLowerName[22:26]:   Core,
	_WorkerName[26:35]:        Crossbell,
	_WorkerLowerName[26:35]:   Crossbell,
	_WorkerName[35:40]:        Curve,
	_WorkerLowerName[35:40]:   Curve,
	_WorkerName[40:43]:        ENS,
	_WorkerLowerName[40:43]:   ENS,
	_WorkerName[43:52]:        Highlight,
	_WorkerLowerName[43:52]:   Highlight,
	_WorkerName[52:58]:        IQWiki,
	_WorkerLowerName[52:58]:   IQWiki,
	_WorkerName[58:67]:        KiwiStand,
	_WorkerLowerName[58:67]:   KiwiStand,
	_WorkerName[67:71]:        Lens,
	_WorkerLowerName[67:71]:   Lens,
	_WorkerName[71:75]:        Lido,
	_WorkerLowerName[71:75]:   Lido,
	_WorkerName[75:84]:        Looksrare,
	_WorkerLowerName[75:84]:   Looksrare,
	_WorkerName[84:91]:        Matters,
	_WorkerLowerName[84:91]:   Matters,
	_WorkerName[91:97]:        Mirror,
	_WorkerLowerName[91:97]:   Mirror,
	_WorkerName[97:103]:       Momoka,
	_WorkerLowerName[97:103]:  Momoka,
	_WorkerName[103:108]:      Oneinch,
	_WorkerLowerName[103:108]: Oneinch,
	_WorkerName[108:115]:      OpenSea,
	_WorkerLowerName[108:115]: OpenSea,
	_WorkerName[115:123]:      Optimism,
	_WorkerLowerName[115:123]: Optimism,
	_WorkerName[123:132]:      Paragraph,
	_WorkerLowerName[123:132]: Paragraph,
	_WorkerName[132:136]:      RSS3,
	_WorkerLowerName[132:136]: RSS3,
	_WorkerName[136:140]:      SAVM,
	_WorkerLowerName[136:140]: SAVM,
	_WorkerName[140:148]:      Stargate,
	_WorkerLowerName[140:148]: Stargate,
	_WorkerName[148:155]:      Uniswap,
	_WorkerLowerName[148:155]: Uniswap,
	_WorkerName[155:158]:      VSL,
	_WorkerLowerName[155:158]: VSL,
}

var _WorkerNames = []string{
	_WorkerName[0:4],
	_WorkerName[4:14],
	_WorkerName[14:22],
	_WorkerName[22:26],
	_WorkerName[26:35],
	_WorkerName[35:40],
	_WorkerName[40:43],
	_WorkerName[43:52],
	_WorkerName[52:58],
	_WorkerName[58:67],
	_WorkerName[67:71],
	_WorkerName[71:75],
	_WorkerName[75:84],
	_WorkerName[84:91],
	_WorkerName[91:97],
	_WorkerName[97:103],
	_WorkerName[103:108],
	_WorkerName[108:115],
	_WorkerName[115:123],
	_WorkerName[123:132],
	_WorkerName[132:136],
	_WorkerName[136:140],
	_WorkerName[140:148],
	_WorkerName[148:155],
	_WorkerName[155:158],
}

// WorkerString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func WorkerString(s string) (Worker, error) {
	if val, ok := _WorkerNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _WorkerNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Worker values", s)
}

// WorkerValues returns all values of the enum
func WorkerValues() []Worker {
	return _WorkerValues
}

// WorkerStrings returns a slice of all String values of the enum
func WorkerStrings() []string {
	strs := make([]string, len(_WorkerNames))
	copy(strs, _WorkerNames)
	return strs
}

// IsAWorker returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Worker) IsAWorker() bool {
	for _, v := range _WorkerValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Worker
func (i Worker) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Worker
func (i *Worker) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Worker should be a string, got %s", data)
	}

	var err error
	*i, err = WorkerString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for Worker
func (i Worker) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Worker
func (i *Worker) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = WorkerString(s)
	return err
}

func (i Worker) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Worker) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of Worker: %[1]T(%[1]v)", value)
	}

	val, err := WorkerString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
