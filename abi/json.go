package abi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/libs4go/errors"
)

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
	module := &moduleImpl{
		encoders:    make(map[string]Encoder),
		funcs:       make(map[string]*funcABI),
		constructor: nil,
	}
	return module, module.parse(json)
}

type funcABI struct {
	selector string  // function selector
	inputs   Encoder // parameter tuple encoder
	outputs  Encoder // return values tuple encoder
}

type moduleImpl struct {
	encoders    map[string]Encoder  // cached encoders
	funcs       map[string]*funcABI // function encoders
	constructor *funcABI            // contructor encoders
}

// parse input json bytes and create memory encoder tree
func (module *moduleImpl) parse(data []byte) error {
	var fields []*JSONField

	err := json.Unmarshal(data, &fields)

	if err != nil {
		return errors.Wrap(err, "parse input json abi errorn")
	}

	for i, field := range fields {
		switch field.Type {
		case JSONTypeFunc:
			abi, err := module.parseFunc(i, field)

			if err != nil {
				return err
			}

			module.funcs[abi.selector] = abi

		case JSONTypeConstructor:
			abi, err := module.parseFunc(i, field)

			if err != nil {
				return err
			}

			module.constructor = abi

		case JSONTypeFallback, JSONTypeReceive:
			// skip parse fallback receive function
		case JSONTypeEvent:
			// TODO: next version support event filter
		default:
			// Skip others

		}
	}

	return nil
}

func (module *moduleImpl) parseFunc(index int, field *JSONField) (*funcABI, error) {
	if field.Type != JSONTypeFunc && field.Type != JSONTypeConstructor {
		return nil, errors.Wrap(ErrJSON, "parse func target(%s) error", field.Type)
	}

	if field.Type == JSONTypeFunc && field.Name == "" {
		return nil, errors.Wrap(ErrJSON, "func name expect,field(%d)", index)
	}

	abi := &funcABI{}

	if field.Type == JSONTypeFunc {
		var buff bytes.Buffer
		buff.WriteString(field.Name)

		buff.WriteString("(")

		var args []string
		for _, arg := range field.Inputs {
			args = append(args, arg.Type)
		}

		buff.WriteString(strings.Join(args, ","))

		buff.WriteString(")")

		abi.selector = hex.EncodeToString(Selector(buff.String()))
	}

	encoder, err := module.parseParams(field.Inputs)

	if err != nil {
		return nil, err
	}

	abi.inputs = encoder

	encoder, err = module.parseParams(field.Outputs)

	if err != nil {
		return nil, err
	}

	abi.outputs = encoder

	return nil, nil
}

func (module *moduleImpl) parseParams(params []*JSONParam) (Encoder, error) {
	var elems []Encoder

	for _, param := range params {
		enc, err := module.parseParam(param)

		if err != nil {
			return nil, err
		}

		elems = append(elems, enc)
	}

	return Tuple(elems...)
}

func (module *moduleImpl) parseParam(param *JSONParam) (Encoder, error) {
	return nil, nil
}
