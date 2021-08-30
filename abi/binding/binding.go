package binding

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/libs4go/errors"
	"github.com/libs4go/ethers/abi"
)

// Some parse regex
var (
	ArrayRegex     = regexp.MustCompile(`([^\[\]]+)|(\[\d*\])`)
	ArrayLenRegex  = regexp.MustCompile(`^\[(\d*)\]$`)
	TupleNameRegex = regexp.MustCompile(`struct ([^\[\]]+)`)
)

type Tuple struct {
	BindingName   string            // golang binding struct name
	BindingFields map[string]string // golang binding struct field
	Encoder       abi.Encoder       // tuple abi encoder
}

type funcABI struct {
	selector []byte      // function selector
	inputs   abi.Encoder // parameter tuple encoder
	outputs  abi.Encoder // return values tuple encoder
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

func (abi *funcABI) Return(data []byte, v interface{}) (uint, error) {
	return abi.outputs.Unmarshal(data, v)
}

type contractImpl struct {
	funcs       map[string]*funcABI // function encoders
	constructor *funcABI            // contructor encoders
}

func (contract *contractImpl) Select(selector string) (abi.Func, bool) {
	f, ok := contract.funcs[selector]

	return f, ok
}

func (contract *contractImpl) parseFunc(index int, field *abi.JSONField, binder Binder) (*funcABI, error) {

	name := strings.Title(field.Name)

	if field.Type == abi.JSONTypeConstructor {
		name = "contructor"
	}

	if field.StateMutability == nil {
		return nil, errors.Wrap(abi.ErrJSON, "func(%s) StateMutability expect", field.Name)
	}

	if field.Type != abi.JSONTypeFunc && field.Type != abi.JSONTypeConstructor {
		return nil, errors.Wrap(abi.ErrJSON, "parse func target(%s) error", field.Type)
	}

	if field.Type == abi.JSONTypeFunc && field.Name == "" {
		return nil, errors.Wrap(abi.ErrJSON, "func name expect,field(%d)", index)
	}

	f := &funcABI{}

	if field.Type == abi.JSONTypeFunc {
		var buff bytes.Buffer
		buff.WriteString(field.Name)

		buff.WriteString("(")

		var args []string
		for _, arg := range field.Inputs {
			args = append(args, arg.Type)
		}

		buff.WriteString(strings.Join(args, ","))

		buff.WriteString(")")

		f.selector = abi.Selector(buff.String())
	}

	encoder, inputs, err := contract.parseParams("inputs", field.Inputs, binder)

	if err != nil {
		return nil, err
	}

	f.inputs = encoder

	encoder, outputs, err := contract.parseParams("outputs", field.Outputs, binder)

	if err != nil {
		return nil, err
	}

	f.outputs = encoder

	binder.Func(name, f.SelectorString(), inputs, outputs, field)

	return f, nil
}

func (contract *contractImpl) parseParams(name string, params []*abi.JSONParam, binder Binder) (abi.Encoder, []abi.Encoder, error) {
	var elems []abi.Encoder

	for _, param := range params {
		enc, err := contract.parseParam(param, binder)

		if err != nil {
			return nil, nil, err
		}

		elems = append(elems, enc)
	}

	t, err := abi.Tuple(name, elems...)

	if err != nil {
		return nil, nil, err
	}

	return t, elems, nil
}

func (contract *contractImpl) parseParam(param *abi.JSONParam, binder Binder) (abi.Encoder, error) {
	encoder, ok := abi.Builtin(param.Type)

	if ok {
		return encoder, nil
	}

	allMatch := ArrayRegex.FindAllString(param.Type, -1)

	if len(allMatch) == 0 {
		return nil, errors.Wrap(abi.ErrJSON, "invalid JSONParam type %s", param.Type)
	}

	if allMatch[0] == "tuple" {
		var err error

		encoder, err = contract.parseTuple(param, binder)

		if err != nil {
			return nil, err
		}

	} else {
		encoder, ok = abi.Builtin(allMatch[0])

		if !ok {
			return nil, errors.Wrap(abi.ErrJSON, "invalid JSONParam type %s", param.Type)
		}
	}

	for _, sub := range allMatch[1:] {
		match := ArrayLenRegex.FindStringSubmatch(sub)

		if len(match) == 0 {
			return nil, errors.Wrap(abi.ErrJSON, "invalid JSONParam type %s, parse submatch %s error", param.Type, sub)
		}

		if match[1] == "" {
			var err error

			encoder, err = abi.Array(encoder)

			if err != nil {
				return nil, err
			}
		} else {
			size, err := strconv.ParseUint(match[1], 10, 64)

			if err != nil {
				return nil, errors.Wrap(abi.ErrJSON, "invalid JSONParam type %s, parse array %s len error", param.Type, sub)
			}

			encoder, err = abi.FixedArray(encoder, uint(size))

			if err != nil {
				return nil, err
			}
		}
	}

	return encoder, nil
}

func (contract *contractImpl) parseTuple(param *abi.JSONParam, binder Binder) (abi.Encoder, error) {

	allMatch := TupleNameRegex.FindStringSubmatch(param.InternalType)

	if len(allMatch) != 2 {
		return nil, errors.Wrap(abi.ErrJSON, "struct InternalType('%s') parse error", param.InternalType)
	}

	tuple, ok := binder.GetTuple(allMatch[1])

	if ok {
		return tuple.Encoder, nil
	}

	var elems []abi.Encoder
	fields := make(map[string]string)

	for _, p := range param.Components {
		elem, err := contract.parseParam(p, binder)

		if err != nil {
			return nil, err
		}

		elems = append(elems, elem)

		fields[strings.Title(p.Name)] = elem.GoTypeName()
	}

	encoder, err := abi.Tuple(allMatch[1], elems...)

	if err != nil {
		return nil, err
	}

	tuple = &Tuple{
		BindingName:   allMatch[1],
		BindingFields: fields,
		Encoder:       encoder,
	}

	binder.RegisterTuple(tuple)

	return encoder, nil
}

// Parse json abi
func Parse(name string, data []byte, binder Binder) (abi.Contract, error) {

	contract := &contractImpl{
		funcs:       make(map[string]*funcABI),
		constructor: nil,
	}

	var fields []*abi.JSONField

	err := json.Unmarshal(data, &fields)

	if err != nil {
		return nil, errors.Wrap(err, "parse input json abi errorn")
	}

	binder.BeginContract(name, data)
	defer binder.EndContract()

	for i, field := range fields {
		switch field.Type {
		case abi.JSONTypeFunc:
			abi, err := contract.parseFunc(i, field, binder)

			if err != nil {
				return nil, err
			}

			contract.funcs[hex.EncodeToString(abi.selector)] = abi

		case abi.JSONTypeConstructor:
			abi, err := contract.parseFunc(i, field, binder)

			if err != nil {
				return nil, err
			}

			contract.constructor = abi

		case abi.JSONTypeFallback, abi.JSONTypeReceive:
			// skip parse fallback receive function
		case abi.JSONTypeEvent:
			// TODO: next version support event filter
		default:
			// Skip others

		}
	}

	return contract, nil
}

func ParseFile(name string, filename string, binder Binder) (abi.Contract, error) {
	buff, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.Wrap(err, "read file: %s error", filename)
	}

	return Parse(name, buff, binder)
}
