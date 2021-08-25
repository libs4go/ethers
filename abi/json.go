package abi

type JSONFieldType string

const (
	JSONTypeFunc        JSONFieldType = "function"
	JSONTypeConstructor JSONFieldType = "constructor"
	JSONTypeReceive     JSONFieldType = "receive"
	JSONTypeFallback    JSONFieldType = "fallback"
	JSONTypeEvent       JSONFieldType = "event"
	JSONTypeError       JSONFieldType = "error"
)

type StateMutability string

const (
	StateMutabilityPure       StateMutability = "pure"
	StateMutabilityView       StateMutability = "view"
	StateMutabilityNonpayable StateMutability = "nonpayable"
	StateMutabilityPayable    StateMutability = "payable"
)

type JSONField struct {
	Type            JSONFieldType    `json:"type"`
	Name            string           `json:"name"`
	Inputs          []*JSONParam     `json:"inputs"`
	Outputs         []*JSONParam     `json:"outputs"`
	StateMutability *StateMutability `json:"stateMutability"`
	Anonymous       *bool            `json:"anonymous"`
}

type JSONParam struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Components []*JSONParam `json:"components"`
	Indexed    *bool        `json:"indexed"`
}

// Module abi parsed module
type Module interface {
}

// Parse abi module from json bytes
func FromJSON(json []byte) (Module, error) {
	module := &moduleImpl{}
	return module, module.parse(json)
}

type moduleImpl struct {
}

// parse input json bytes and create memory encoder tree
func (module *moduleImpl) parse(json []byte) error {
	return nil
}
