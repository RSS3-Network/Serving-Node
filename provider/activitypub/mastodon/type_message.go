// Code generated by "enumer --values --type=MessageType --output type_message.go --trimprefix=MessageType"; DO NOT EDIT.

package mastodon

import (
	"fmt"
	"strings"
)

const _MessageTypeName = "NoneCreateAnnounceLike"

var _MessageTypeIndex = [...]uint8{0, 4, 10, 18, 22}

const _MessageTypeLowerName = "nonecreateannouncelike"

func (i MessageType) String() string {
	if i < 0 || i >= MessageType(len(_MessageTypeIndex)-1) {
		return fmt.Sprintf("MessageType(%d)", i)
	}
	return _MessageTypeName[_MessageTypeIndex[i]:_MessageTypeIndex[i+1]]
}

func (MessageType) Values() []string {
	return MessageTypeStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _MessageTypeNoOp() {
	var x [1]struct{}
	_ = x[MessageTypeNone-(0)]
	_ = x[MessageTypeCreate-(1)]
	_ = x[MessageTypeAnnounce-(2)]
	_ = x[MessageTypeLike-(3)]
}

var _MessageTypeValues = []MessageType{MessageTypeNone, MessageTypeCreate, MessageTypeAnnounce, MessageTypeLike}

var _MessageTypeNameToValueMap = map[string]MessageType{
	_MessageTypeName[0:4]:        MessageTypeNone,
	_MessageTypeLowerName[0:4]:   MessageTypeNone,
	_MessageTypeName[4:10]:       MessageTypeCreate,
	_MessageTypeLowerName[4:10]:  MessageTypeCreate,
	_MessageTypeName[10:18]:      MessageTypeAnnounce,
	_MessageTypeLowerName[10:18]: MessageTypeAnnounce,
	_MessageTypeName[18:22]:      MessageTypeLike,
	_MessageTypeLowerName[18:22]: MessageTypeLike,
}

var _MessageTypeNames = []string{
	_MessageTypeName[0:4],
	_MessageTypeName[4:10],
	_MessageTypeName[10:18],
	_MessageTypeName[18:22],
}

// MessageTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func MessageTypeString(s string) (MessageType, error) {
	if val, ok := _MessageTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _MessageTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to MessageType values", s)
}

// MessageTypeValues returns all values of the enum
func MessageTypeValues() []MessageType {
	return _MessageTypeValues
}

// MessageTypeStrings returns a slice of all String values of the enum
func MessageTypeStrings() []string {
	strs := make([]string, len(_MessageTypeNames))
	copy(strs, _MessageTypeNames)
	return strs
}

// IsAMessageType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i MessageType) IsAMessageType() bool {
	for _, v := range _MessageTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
