// Code generated by "enumer --values --type=SocialType --transform=snake --trimprefix=TypeSocial --output type_social_string.go --json --sql"; DO NOT EDIT.

package filter

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _SocialTypeName = "postrevisecommentshareprofile"

var _SocialTypeIndex = [...]uint8{0, 4, 10, 17, 22, 29}

const _SocialTypeLowerName = "postrevisecommentshareprofile"

func (i SocialType) String() string {
	i -= 1
	if i >= SocialType(len(_SocialTypeIndex)-1) {
		return fmt.Sprintf("SocialType(%d)", i+1)
	}
	return _SocialTypeName[_SocialTypeIndex[i]:_SocialTypeIndex[i+1]]
}

func (SocialType) Values() []string {
	return SocialTypeStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _SocialTypeNoOp() {
	var x [1]struct{}
	_ = x[TypeSocialPost-(1)]
	_ = x[TypeSocialRevise-(2)]
	_ = x[TypeSocialComment-(3)]
	_ = x[TypeSocialShare-(4)]
	_ = x[TypeSocialProfile-(5)]
}

var _SocialTypeValues = []SocialType{TypeSocialPost, TypeSocialRevise, TypeSocialComment, TypeSocialShare, TypeSocialProfile}

var _SocialTypeNameToValueMap = map[string]SocialType{
	_SocialTypeName[0:4]:        TypeSocialPost,
	_SocialTypeLowerName[0:4]:   TypeSocialPost,
	_SocialTypeName[4:10]:       TypeSocialRevise,
	_SocialTypeLowerName[4:10]:  TypeSocialRevise,
	_SocialTypeName[10:17]:      TypeSocialComment,
	_SocialTypeLowerName[10:17]: TypeSocialComment,
	_SocialTypeName[17:22]:      TypeSocialShare,
	_SocialTypeLowerName[17:22]: TypeSocialShare,
	_SocialTypeName[22:29]:      TypeSocialProfile,
	_SocialTypeLowerName[22:29]: TypeSocialProfile,
}

var _SocialTypeNames = []string{
	_SocialTypeName[0:4],
	_SocialTypeName[4:10],
	_SocialTypeName[10:17],
	_SocialTypeName[17:22],
	_SocialTypeName[22:29],
}

// SocialTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SocialTypeString(s string) (SocialType, error) {
	if val, ok := _SocialTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _SocialTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to SocialType values", s)
}

// SocialTypeValues returns all values of the enum
func SocialTypeValues() []SocialType {
	return _SocialTypeValues
}

// SocialTypeStrings returns a slice of all String values of the enum
func SocialTypeStrings() []string {
	strs := make([]string, len(_SocialTypeNames))
	copy(strs, _SocialTypeNames)
	return strs
}

// IsASocialType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i SocialType) IsASocialType() bool {
	for _, v := range _SocialTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for SocialType
func (i SocialType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for SocialType
func (i *SocialType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("SocialType should be a string, got %s", data)
	}

	var err error
	*i, err = SocialTypeString(s)
	return err
}

func (i SocialType) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *SocialType) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of SocialType: %[1]T(%[1]v)", value)
	}

	val, err := SocialTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
