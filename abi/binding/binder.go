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
	Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder, jsondata *abi.JSONField)
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

func (impl *symbolsImpl) Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder, jsondata *abi.JSONField) {

}

func (impl *symbolsImpl) EndContract() {

}

type Contract struct {
	Name          string
	ABI           string
	Funcs         []*Func
	Contructor    *Func
	overrideFuncs map[string]int
}

type Func struct {
	ReadOnly       bool
	Name           string
	Selector       string
	Inputs         []string
	Outputs        []string
	GoInputParams  string
	GoOutputParams string
	GoInputArgs    string
	GoOutputArgs   string
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
		Name:          name,
		ABI:           hex.EncodeToString(abi),
		overrideFuncs: make(map[string]int),
	})
}

func (impl *Generator) Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder, jsondata *abi.JSONField) {

	readOnly := true

	state := *jsondata.StateMutability

	contructor := jsondata.Type == abi.JSONTypeConstructor

	if state == abi.StateMutabilityNonpayable || state == abi.StateMutabilityPayable {
		readOnly = false
	}

	if len(impl.contracts) == 0 {
		return
	}

	c := impl.contracts[len(impl.contracts)-1]

	if counter, ok := c.overrideFuncs[name]; ok {
		c.overrideFuncs[name] = counter + 1
		name = fmt.Sprintf("%s%d", name, counter)
	} else {
		c.overrideFuncs[name] = 1
	}

	var is []string
	var inputParams []string
	var inputArgs []string

	for i, p := range inputs {
		is = append(is, p.GoTypeName())

		name := jsondata.Inputs[i].Name

		if name == "" {
			name = fmt.Sprintf("param%d", i)
		}

		inputParams = append(inputParams, fmt.Sprintf("%s %s", name, p.GoTypeName()))

		inputArgs = append(inputArgs, name)
	}

	if !readOnly {
		inputParams = append(inputParams, "ops ...abi.Op")
	}

	var os []string
	var outputParams []string
	var goOutputArgs []string

	for i, p := range outputs {

		if contructor {
			continue
		}

		is = append(os, p.GoTypeName())

		if readOnly {
			name := jsondata.Outputs[i].Name

			if name == "" {
				name = fmt.Sprintf("ret%d", i)
			}
			outputParams = append(outputParams, fmt.Sprintf("%s %s", name, p.GoTypeName()))
			goOutputArgs = append(goOutputArgs, fmt.Sprintf("&%s", name))
		}

	}

	if !readOnly {
		outputParams = append(outputParams, "ret0 abi.Transaction")
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
			GoInputArgs:    strings.Join(inputArgs, ", "),
			GoOutputArgs:   strings.Join(goOutputArgs, ", "),
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

func (impl *Generator) calcImports(buff bytes.Buffer) []string {
	content := buff.String()

	var imports []string

	if strings.Contains(content, "*big.Int") {
		imports = append(imports, "math/big")
	}

	if strings.Contains(content, "address.Address") {
		imports = append(imports, "github.com/libs4go/ethers/address")
	}

	return imports
}

type headerModel struct {
	Name    string
	Imports []string
}

func (impl *Generator) Write(packageName string, writer io.Writer) error {

	var buff bytes.Buffer

	err := tupleTmpl.Execute(&buff, impl.tuples)

	if err != nil {
		return err
	}

	err = contractTmpl.Execute(&buff, impl.contracts)

	if err != nil {
		return err
	}

	hm := &headerModel{
		Name:    packageName,
		Imports: impl.calcImports(buff),
	}

	var headerBuff bytes.Buffer

	err = headerTmpl.Execute(&headerBuff, hm)

	if err != nil {
		return err
	}

	_, err = writer.Write(headerBuff.Bytes())

	if err != nil {
		return err
	}

	_, err = writer.Write(buff.Bytes())

	if err != nil {
		return err
	}

	return err
}

func NewGen() *Generator {
	return &Generator{
		tuples: make(map[string]*Tuple),
	}
}

var headerTmplText = `
package {{.Name}}

import (
	"context"
	"encoding/hex"

	"github.com/libs4go/errors"
	"github.com/libs4go/ethers/abi"
	"github.com/libs4go/ethers/abi/binding"
	"github.com/libs4go/ethers/client"
	"github.com/libs4go/ethers/signer"
	

{{range $_, $import := .Imports}}
	"{{$import}}"
{{end}}

)
`

var headerTmpl = template.Must(template.New("Gen").Parse(headerTmplText))

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
	{{$field.Name}}(ctx context.Context, {{$field.GoInputParams}})({{$field.GoOutputParams}})
	{{end}}
}

type impl{{$element.Name}} struct {
	contract abi.Contract
	client client.Provider
	signer signer.Signer
	recipient string
}

{{range $_, $field := $element.Funcs}}
func (impl *impl{{$element.Name}}) {{$field.Name}}(ctx context.Context, {{$field.GoInputParams}})({{$field.GoOutputParams}}) {
	f, ok :=  abi.TryGetFunc(impl.contract, "{{$field.Selector}}")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func {{$field.Name}} not found")
		return
	}

	var buff []byte

	buff, err = f.Call({{$field.GoInputArgs}})

	if err != nil {
		return
	}

	{{if $field.ReadOnly}}
	callSite := &client.CallSite {
		To: impl.recipient,
		Data: hex.EncodeToString(buff),
	}
	
	var ret string
	
	ret,err = impl.client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff,[]interface{}{ {{$field.GoOutputArgs}} })

	return

	{{else}}
	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

	{{end}}
}
{{end}}

{{end}}
`

var contractTmpl = template.Must(template.New("Gen").Parse(contractTmplText))
