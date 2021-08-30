package binding

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/libs4go/ethers/abi"
)

// ABI binder context
type Binder interface {
	RegisterTuple(tuple *Tuple)
	GetTuple(name string) (*Tuple, bool)
	BeginContract(name string, abi []byte)
	Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder, contructor bool, state abi.StateMutability)
	EndContract()
}

type symbolsImpl struct {
	tuples map[string]*Tuple
}

// NewSymbols create binder binder
func NewSymbols() Binder {
	return &symbolsImpl{
		tuples: make(map[string]*Tuple),
	}
}

func (impl *symbolsImpl) RegisterTuple(tuple *Tuple) {
	impl.tuples[tuple.BindingName] = tuple
}

func (impl *symbolsImpl) GetTuple(name string) (*Tuple, bool) {
	t, ok := impl.tuples[name]

	return t, ok
}

func (impl *symbolsImpl) BeginContract(name string, abi []byte) {

}

func (impl *symbolsImpl) Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder, contructor bool, state abi.StateMutability) {

}

func (impl *symbolsImpl) EndContract() {

}

type Contract struct {
	Name       string
	ABI        string
	Funcs      []*Func
	Contructor *Func
}

type Func struct {
	ReadOnly       bool
	Name           string
	Selector       string
	Inputs         []string
	Outputs        []string
	GoInputParams  string
	GoOutputParams string
}

type Generator struct {
	tuples    map[string]*Tuple
	contracts []*Contract
}

func (impl *Generator) RegisterTuple(tuple *Tuple) {
	impl.tuples[tuple.BindingName] = tuple
}

func (impl *Generator) GetTuple(name string) (*Tuple, bool) {
	t, ok := impl.tuples[name]

	return t, ok
}

func (impl *Generator) BeginContract(name string, abi []byte) {
	impl.contracts = append(impl.contracts, &Contract{
		Name: name,
		ABI:  hex.EncodeToString(abi),
	})
}

func (impl *Generator) Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder, contructor bool, state abi.StateMutability) {

	readOnly := true

	if state == abi.StateMutabilityNonpayable || state == abi.StateMutabilityPayable {
		readOnly = false
	}

	if len(impl.contracts) == 0 {
		return
	}

	c := impl.contracts[len(impl.contracts)-1]

	var is []string
	var inputParams []string

	for i, p := range inputs {
		is = append(is, p.GoTypeName())

		if name == "Hello" {
			println(p.GoTypeName())
		}

		inputParams = append(inputParams, fmt.Sprintf("param%d %s", i, p.GoTypeName()))
	}

	var os []string
	var outputParams []string

	for i, p := range outputs {

		if contructor {
			continue
		}

		is = append(os, p.GoTypeName())

		if readOnly {
			outputParams = append(outputParams, fmt.Sprintf("ret%d %s", i, p.GoTypeName()))
		}

	}

	if !readOnly {
		outputParams = append(outputParams, "ret0 *abi.Transaction")
	}

	if !contructor {
		outputParams = append(outputParams, "err error")

		c.Funcs = append(c.Funcs, &Func{
			ReadOnly:       readOnly,
			Name:           name,
			Selector:       selector,
			Inputs:         is,
			Outputs:        os,
			GoInputParams:  strings.Join(inputParams, ", "),
			GoOutputParams: strings.Join(outputParams, ", "),
		})
	} else {
		c.Contructor = &Func{
			Name:           "New",
			Selector:       selector,
			Inputs:         is,
			Outputs:        os,
			GoInputParams:  strings.Join(inputParams, ", "),
			GoOutputParams: strings.Join(outputParams, ", "),
		}
	}

}

func (impl *Generator) EndContract() {

}

func (impl *Generator) Write(writer io.Writer) error {

	var buff bytes.Buffer

	err := tupleTmpl.Execute(&buff, impl.tuples)

	if err != nil {
		return err
	}

	err = contractTmpl.Execute(&buff, impl.contracts)

	if err != nil {
		return err
	}

	_, err = writer.Write(buff.Bytes())

	return err
}

func NewGen() *Generator {
	return &Generator{
		tuples: make(map[string]*Tuple),
	}
}

var tupleTmplText = `
{{range $key, $element := .}}
// Generated tuple "{{$key}}" stub code , do not modify manually
type {{$key}} struct {
	{{range $name, $field := $element.BindingFields}}
	{{$name}} {{$field}}
	{{end}}
}

{{end}}
`

var tupleTmpl = template.Must(template.New("Gen").Parse(tupleTmplText))

var contractTmplText = `
{{range $index, $element := .}}
type {{$element.Name}} interface {
	{{range $_, $field := $element.Funcs}}
	{{$field.Name}}({{$field.GoInputParams}})({{$field.GoOutputParams}})
	{{end}}
}

{{end}}
`

var contractTmpl = template.Must(template.New("Gen").Parse(contractTmplText))
