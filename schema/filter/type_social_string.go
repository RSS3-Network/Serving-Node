// Code generated by "enumer --values --type=SocialType --transform=snake --trimprefix=TypeSocial --output type_social_string.go --json --sql"; DO NOT EDIT.

package filter

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _SocialTypeName = "postcommentsharemintprofileproxyrevisedeletereward"

var _SocialTypeIndex = [...]uint8{0, 4, 11, 16, 20, 27, 32, 38, 44, 50}

const _SocialTypeLowerName = "postcommentsharemintprofileproxyrevisedeletereward"

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
	_ = x[TypeSocialComment-(2)]
	_ = x[TypeSocialShare-(3)]
	_ = x[TypeSocialMint-(4)]
	_ = x[TypeSocialProfile-(5)]
	_ = x[TypeSocialProxy-(6)]
	_ = x[TypeSocialRevise-(7)]
	_ = x[TypeSocialDelete-(8)]
	_ = x[TypeSocialReward-(9)]
}

var _SocialTypeValues = []SocialType{TypeSocialPost, TypeSocialComment, TypeSocialShare, TypeSocialMint, TypeSocialProfile, TypeSocialProxy, TypeSocialRevise, TypeSocialDelete, TypeSocialReward}

var _SocialTypeNameToValueMap = map[string]SocialType{
	_SocialTypeName[0:4]:        TypeSocialPost,
	_SocialTypeLowerName[0:4]:   TypeSocialPost,
	_SocialTypeName[4:11]:       TypeSocialComment,
	_SocialTypeLowerName[4:11]:  TypeSocialComment,
	_SocialTypeName[11:16]:      TypeSocialShare,
	_SocialTypeLowerName[11:16]: TypeSocialShare,
	_SocialTypeName[16:20]:      TypeSocialMint,
	_SocialTypeLowerName[16:20]: TypeSocialMint,
	_SocialTypeName[20:27]:      TypeSocialProfile,
	_SocialTypeLowerName[20:27]: TypeSocialProfile,
	_SocialTypeName[27:32]:      TypeSocialProxy,
	_SocialTypeLowerName[27:32]: TypeSocialProxy,
	_SocialTypeName[32:38]:      TypeSocialRevise,
	_SocialTypeLowerName[32:38]: TypeSocialRevise,
	_SocialTypeName[38:44]:      TypeSocialDelete,
	_SocialTypeLowerName[38:44]: TypeSocialDelete,
	_SocialTypeName[44:50]:      TypeSocialReward,
	_SocialTypeLowerName[44:50]: TypeSocialReward,
}

var _SocialTypeNames = []string{
	_SocialTypeName[0:4],
	_SocialTypeName[4:11],
	_SocialTypeName[11:16],
	_SocialTypeName[16:20],
	_SocialTypeName[20:27],
	_SocialTypeName[27:32],
	_SocialTypeName[32:38],
	_SocialTypeName[38:44],
	_SocialTypeName[44:50],
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
