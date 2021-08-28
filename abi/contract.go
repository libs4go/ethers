package abi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/libs4go/errors"
)

var basicTypeEncoders map[string]Encoder
var basicTypeOnce sync.Once

func ensure(encoder Encoder, err error) Encoder {
	if err != nil {
		panic(err)
	}

	return encoder
}

func initBasicTypeEncoders() {
	basicTypeEncoders = make(map[string]Encoder)

	var encoders []Encoder

	for i := uint(0); i < 32; i++ {
		encoders = append(encoders, ensure(Integer(false, (i+1)*8)))
		encoders = append(encoders, ensure(Integer(true, (i+1)*8)))
		encoders = append(encoders, ensure(FixedBytes((i + 1))))
	}

	encoders = append(encoders, ensure(Bool()))

	encoders = append(encoders, ensure(String()))

	encoders = append(encoders, ensure(Bytes()))

	for _, encoder := range encoders {
		basicTypeEncoders[encoder.String()] = encoder
	}

	basicTypeEncoders["address"] = basicTypeEncoders["uint160"]
	basicTypeEncoders["uint"] = basicTypeEncoders["uint256"]
	basicTypeEncoders["int"] = basicTypeEncoders["int256"]
}

func basicTypes() map[string]Encoder {
	basicTypeOnce.Do(initBasicTypeEncoders)

	return basicTypeEncoders
}

// Module abi parsed contract
type Contract interface {
	// Get func by signature
	Func(signature string) (Func, bool)
}

type Func interface {
	Selector() []byte
	SelectorString() string
	Call(params ...interface{}) ([]byte, error)
	Returns(data []byte, v interface{}) (uint, error)
}

// Parse abi contract from json bytes
func Parse(json []byte) (Contract, error) {
	contract := &contractImpl{
		tuples:      make(map[string]Encoder),
		funcs:       make(map[string]*funcABI),
		constructor: nil,
	}
	return contract, contract.parse(json)
}

func ParseFile(filename string) (Contract, error) {
	buff, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.Wrap(err, "read file: %s error", filename)
	}

	return Parse(buff)
}

type funcABI struct {
	selector []byte  // function selector
	inputs   Encoder // parameter tuple encoder
	outputs  Encoder // return values tuple encoder
}

func (abi *funcABI) Selector() []byte {
	return abi.selector
}

func (abi *funcABI) SelectorString() string {
	return hex.EncodeToString(abi.selector)
}

func (abi *funcABI) Call(params ...interface{}) ([]byte, error) {
	buff, err := abi.inputs.Marshal(params)

	if err != nil {
		return nil, err
	}

	return append(abi.selector, buff...), nil
}

func (abi *funcABI) Returns(data []byte, v interface{}) (uint, error) {
	return abi.outputs.Unmarshal(data, v)
}

type contractImpl struct {
	tuples      map[string]Encoder  // cached encoders
	funcs       map[string]*funcABI // function encoders
	constructor *funcABI            // contructor encoders
}

func (contract *contractImpl) Func(signature string) (Func, bool) {
	f, ok := contract.funcs[hex.EncodeToString(Selector(signature))]

	return f, ok
}

// parse input json bytes and create memory encoder tree
func (contract *contractImpl) parse(data []byte) error {
	var fields []*JSONField

	err := json.Unmarshal(data, &fields)

	if err != nil {
		return errors.Wrap(err, "parse input json abi errorn")
	}

	for i, field := range fields {
		switch field.Type {
		case JSONTypeFunc:
			abi, err := contract.parseFunc(i, field)

			if err != nil {
				return err
			}

			contract.funcs[hex.EncodeToString(abi.selector)] = abi

		case JSONTypeConstructor:
			abi, err := contract.parseFunc(i, field)

			if err != nil {
				return err
			}

			contract.constructor = abi

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

func (contract *contractImpl) parseFunc(index int, field *JSONField) (*funcABI, error) {
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

		abi.selector = Selector(buff.String())
	}

	encoder, err := contract.parseParams("inputs", field.Inputs)

	if err != nil {
		return nil, err
	}

	abi.inputs = encoder

	encoder, err = contract.parseParams("outputs", field.Outputs)

	if err != nil {
		return nil, err
	}

	abi.outputs = encoder

	return abi, nil
}

func (contract *contractImpl) parseParams(name string, params []*JSONParam) (Encoder, error) {
	var elems []Encoder

	for _, param := range params {
		enc, err := contract.parseParam(param)

		if err != nil {
			return nil, err
		}

		elems = append(elems, enc)
	}

	return Tuple(name, elems...)
}

var arrayTypeRegex = regexp.MustCompile(`([^\[\]]+)|(\[\d*\])`)
var arrayTypeLenRegex = regexp.MustCompile(`^\[(\d*)\]$`)
var tupleNameRegex = regexp.MustCompile(`struct ([^\[\]]+)`)

func (contract *contractImpl) parseParam(param *JSONParam) (Encoder, error) {
	encoder, ok := basicTypes()[param.Type]

	if ok {
		return encoder, nil
	}

	allMatch := arrayTypeRegex.FindAllString(param.Type, -1)

	if len(allMatch) == 0 {
		return nil, errors.Wrap(ErrJSON, "invalid JSONParam type %s", param.Type)
	}

	if allMatch[0] == "tuple" {
		var err error

		encoder, err = contract.parseTuple(param)

		if err != nil {
			return nil, err
		}

	} else {
		encoder, ok = basicTypes()[allMatch[0]]

		if !ok {
			return nil, errors.Wrap(ErrJSON, "invalid JSONParam type %s", param.Type)
		}
	}

	for _, sub := range allMatch[1:] {
		match := arrayTypeLenRegex.FindStringSubmatch(sub)

		if len(match) == 0 {
			return nil, errors.Wrap(ErrJSON, "invalid JSONParam type %s, parse submatch %s error", param.Type, sub)
		}

		if match[1] == "" {
			var err error

			encoder, err = Array(encoder)

			if err != nil {
				return nil, err
			}
		} else {
			size, err := strconv.ParseUint(match[1], 10, 64)

			if err != nil {
				return nil, errors.Wrap(ErrJSON, "invalid JSONParam type %s, parse array %s len error", param.Type, sub)
			}

			encoder, err = FixedArray(encoder, uint(size))

			if err != nil {
				return nil, err
			}
		}
	}

	return encoder, nil
}

func (contract *contractImpl) parseTuple(param *JSONParam) (Encoder, error) {

	allMatch := tupleNameRegex.FindStringSubmatch(param.InternalType)

	if len(allMatch) != 2 {
		return nil, errors.Wrap(ErrJSON, "struct InternalType('%s') parse error", param.InternalType)
	}

	encoder, ok := contract.tuples[allMatch[1]]

	if ok {
		return encoder, nil
	}

	var elems []Encoder

	for _, p := range param.Components {
		elem, err := contract.parseParam(p)

		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	encoder, err := Tuple(allMatch[1], elems...)

	if err != nil {
		return nil, err
	}

	contract.tuples[allMatch[1]] = encoder

	return encoder, nil
}
