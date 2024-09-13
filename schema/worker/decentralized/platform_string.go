// Code generated by "enumer --values --type=Platform --linecomment --output platform_string.go --json --yaml --sql"; DO NOT EDIT.

package decentralized

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _PlatformName = "Unknown1inchAAVEAavegotchiArbitrumBaseBendDAOCowCrossbellCurveENSFarcasterHighlightIQWikiKiwiStandLensLidoLinearLiNEARLooksRareMattersMirrorNounsOpenSeaOptimismParagraphParaswapRSS3SAVMStargateUniswapVSL"

var _PlatformIndex = [...]uint8{0, 7, 12, 16, 26, 34, 38, 45, 48, 57, 62, 65, 74, 83, 89, 98, 102, 106, 112, 118, 127, 134, 140, 145, 152, 160, 169, 177, 181, 185, 193, 200, 203}

const _PlatformLowerName = "unknown1inchaaveaavegotchiarbitrumbasebenddaocowcrossbellcurveensfarcasterhighlightiqwikikiwistandlenslidolinearlinearlooksraremattersmirrornounsopenseaoptimismparagraphparaswaprss3savmstargateuniswapvsl"

func (i Platform) String() string {
	if i >= Platform(len(_PlatformIndex)-1) {
		return fmt.Sprintf("Platform(%d)", i)
	}
	return _PlatformName[_PlatformIndex[i]:_PlatformIndex[i+1]]
}

func (Platform) Values() []string {
	return PlatformStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _PlatformNoOp() {
	var x [1]struct{}
	_ = x[PlatformUnknown-(0)]
	_ = x[Platform1Inch-(1)]
	_ = x[PlatformAAVE-(2)]
	_ = x[PlatformAavegotchi-(3)]
	_ = x[PlatformArbitrum-(4)]
	_ = x[PlatformBase-(5)]
	_ = x[PlatformBendDAO-(6)]
	_ = x[PlatformCow-(7)]
	_ = x[PlatformCrossbell-(8)]
	_ = x[PlatformCurve-(9)]
	_ = x[PlatformENS-(10)]
	_ = x[PlatformFarcaster-(11)]
	_ = x[PlatformHighlight-(12)]
	_ = x[PlatformIQWiki-(13)]
	_ = x[PlatformKiwiStand-(14)]
	_ = x[PlatformLens-(15)]
	_ = x[PlatformLido-(16)]
	_ = x[PlatformLinea-(17)]
	_ = x[PlatformLiNEAR-(18)]
	_ = x[PlatformLooksRare-(19)]
	_ = x[PlatformMatters-(20)]
	_ = x[PlatformMirror-(21)]
	_ = x[PlatformNouns-(22)]
	_ = x[PlatformOpenSea-(23)]
	_ = x[PlatformOptimism-(24)]
	_ = x[PlatformParagraph-(25)]
	_ = x[PlatformParaswap-(26)]
	_ = x[PlatformRSS3-(27)]
	_ = x[PlatformSAVM-(28)]
	_ = x[PlatformStargate-(29)]
	_ = x[PlatformUniswap-(30)]
	_ = x[PlatformVSL-(31)]
}

var _PlatformValues = []Platform{PlatformUnknown, Platform1Inch, PlatformAAVE, PlatformAavegotchi, PlatformArbitrum, PlatformBase, PlatformBendDAO, PlatformCow, PlatformCrossbell, PlatformCurve, PlatformENS, PlatformFarcaster, PlatformHighlight, PlatformIQWiki, PlatformKiwiStand, PlatformLens, PlatformLido, PlatformLinea, PlatformLiNEAR, PlatformLooksRare, PlatformMatters, PlatformMirror, PlatformNouns, PlatformOpenSea, PlatformOptimism, PlatformParagraph, PlatformParaswap, PlatformRSS3, PlatformSAVM, PlatformStargate, PlatformUniswap, PlatformVSL}

var _PlatformNameToValueMap = map[string]Platform{
	_PlatformName[0:7]:          PlatformUnknown,
	_PlatformLowerName[0:7]:     PlatformUnknown,
	_PlatformName[7:12]:         Platform1Inch,
	_PlatformLowerName[7:12]:    Platform1Inch,
	_PlatformName[12:16]:        PlatformAAVE,
	_PlatformLowerName[12:16]:   PlatformAAVE,
	_PlatformName[16:26]:        PlatformAavegotchi,
	_PlatformLowerName[16:26]:   PlatformAavegotchi,
	_PlatformName[26:34]:        PlatformArbitrum,
	_PlatformLowerName[26:34]:   PlatformArbitrum,
	_PlatformName[34:38]:        PlatformBase,
	_PlatformLowerName[34:38]:   PlatformBase,
	_PlatformName[38:45]:        PlatformBendDAO,
	_PlatformLowerName[38:45]:   PlatformBendDAO,
	_PlatformName[45:48]:        PlatformCow,
	_PlatformLowerName[45:48]:   PlatformCow,
	_PlatformName[48:57]:        PlatformCrossbell,
	_PlatformLowerName[48:57]:   PlatformCrossbell,
	_PlatformName[57:62]:        PlatformCurve,
	_PlatformLowerName[57:62]:   PlatformCurve,
	_PlatformName[62:65]:        PlatformENS,
	_PlatformLowerName[62:65]:   PlatformENS,
	_PlatformName[65:74]:        PlatformFarcaster,
	_PlatformLowerName[65:74]:   PlatformFarcaster,
	_PlatformName[74:83]:        PlatformHighlight,
	_PlatformLowerName[74:83]:   PlatformHighlight,
	_PlatformName[83:89]:        PlatformIQWiki,
	_PlatformLowerName[83:89]:   PlatformIQWiki,
	_PlatformName[89:98]:        PlatformKiwiStand,
	_PlatformLowerName[89:98]:   PlatformKiwiStand,
	_PlatformName[98:102]:       PlatformLens,
	_PlatformLowerName[98:102]:  PlatformLens,
	_PlatformName[102:106]:      PlatformLido,
	_PlatformLowerName[102:106]: PlatformLido,
	_PlatformName[106:112]:      PlatformLinea,
	_PlatformLowerName[106:112]: PlatformLinea,
	_PlatformName[112:118]:      PlatformLiNEAR,
	_PlatformLowerName[112:118]: PlatformLiNEAR,
	_PlatformName[118:127]:      PlatformLooksRare,
	_PlatformLowerName[118:127]: PlatformLooksRare,
	_PlatformName[127:134]:      PlatformMatters,
	_PlatformLowerName[127:134]: PlatformMatters,
	_PlatformName[134:140]:      PlatformMirror,
	_PlatformLowerName[134:140]: PlatformMirror,
	_PlatformName[140:145]:      PlatformNouns,
	_PlatformLowerName[140:145]: PlatformNouns,
	_PlatformName[145:152]:      PlatformOpenSea,
	_PlatformLowerName[145:152]: PlatformOpenSea,
	_PlatformName[152:160]:      PlatformOptimism,
	_PlatformLowerName[152:160]: PlatformOptimism,
	_PlatformName[160:169]:      PlatformParagraph,
	_PlatformLowerName[160:169]: PlatformParagraph,
	_PlatformName[169:177]:      PlatformParaswap,
	_PlatformLowerName[169:177]: PlatformParaswap,
	_PlatformName[177:181]:      PlatformRSS3,
	_PlatformLowerName[177:181]: PlatformRSS3,
	_PlatformName[181:185]:      PlatformSAVM,
	_PlatformLowerName[181:185]: PlatformSAVM,
	_PlatformName[185:193]:      PlatformStargate,
	_PlatformLowerName[185:193]: PlatformStargate,
	_PlatformName[193:200]:      PlatformUniswap,
	_PlatformLowerName[193:200]: PlatformUniswap,
	_PlatformName[200:203]:      PlatformVSL,
	_PlatformLowerName[200:203]: PlatformVSL,
}

var _PlatformNames = []string{
	_PlatformName[0:7],
	_PlatformName[7:12],
	_PlatformName[12:16],
	_PlatformName[16:26],
	_PlatformName[26:34],
	_PlatformName[34:38],
	_PlatformName[38:45],
	_PlatformName[45:48],
	_PlatformName[48:57],
	_PlatformName[57:62],
	_PlatformName[62:65],
	_PlatformName[65:74],
	_PlatformName[74:83],
	_PlatformName[83:89],
	_PlatformName[89:98],
	_PlatformName[98:102],
	_PlatformName[102:106],
	_PlatformName[106:112],
	_PlatformName[112:118],
	_PlatformName[118:127],
	_PlatformName[127:134],
	_PlatformName[134:140],
	_PlatformName[140:145],
	_PlatformName[145:152],
	_PlatformName[152:160],
	_PlatformName[160:169],
	_PlatformName[169:177],
	_PlatformName[177:181],
	_PlatformName[181:185],
	_PlatformName[185:193],
	_PlatformName[193:200],
	_PlatformName[200:203],
}

// PlatformString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PlatformString(s string) (Platform, error) {
	if val, ok := _PlatformNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _PlatformNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Platform values", s)
}

// PlatformValues returns all values of the enum
func PlatformValues() []Platform {
	return _PlatformValues
}

// PlatformStrings returns a slice of all String values of the enum
func PlatformStrings() []string {
	strs := make([]string, len(_PlatformNames))
	copy(strs, _PlatformNames)
	return strs
}

// IsAPlatform returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Platform) IsAPlatform() bool {
	for _, v := range _PlatformValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Platform
func (i Platform) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Platform
func (i *Platform) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Platform should be a string, got %s", data)
	}

	var err error
	*i, err = PlatformString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for Platform
func (i Platform) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Platform
func (i *Platform) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = PlatformString(s)
	return err
}

func (i Platform) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Platform) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of Platform: %[1]T(%[1]v)", value)
	}

	val, err := PlatformString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
