package dynamic

import (
	"context"
	"encoding/json"

	"github.com/libs4go/errors"
	"github.com/libs4go/ethers/abi"
	"github.com/libs4go/ethers/client"
	"github.com/libs4go/ethers/signer"
)

// Contract dynamic evm abi callsite
type Contract interface {
	Name() string
	// Connect new rpc client and signer
	Connect(address string, provider client.Provider)
	// Contract read operation
	Call(ctx context.Context, method string, result interface{}, args ...interface{}) error
	// Contract write operation
	Exec(ctx context.Context, signer signer.Signer, method string, args ...interface{}) error
}

type paramSite struct {
}

type errorSite struct {
	Name   string       `json:"name"`
	Inputs []*paramSite `json:"inputs"`
}

type eventSite struct {
	Name      string       `json:"name"`
	Inputs    []*paramSite `json:"inputs"`
	Anonymous *bool        `json:"anonymous"`
}

type execSite struct {
	Name    string
	Inputs  []*paramSite
	Outputs []*paramSite
	Payable bool
}

type callSite struct {
	Name    string
	Inputs  []*paramSite
	Outputs []*paramSite
	Pure    bool
}

type contractImpl struct {
	address     string
	provider    client.Provider
	name        string
	constructor *execSite
	fallback    *execSite
	receive     *execSite
	execs       map[string]*execSite
	calls       map[string]*callSite
	events      map[string]*eventSite
	errors      map[string]*errorSite
}

func toExecSite(field *abi.JSONField) *execSite {
	return nil
}

func toCallSite(field *abi.JSONField) *callSite {
	return nil
}

func toErrorSite(field *abi.JSONField) *errorSite {
	return nil
}

func toEventSite(field *abi.JSONField) *eventSite {
	return nil
}

// Load JSON abi and set contract name
func Load(name string, abiJSON []byte) (Contract, error) {
	var fields []abi.JSONField

	err := json.Unmarshal(abiJSON, &fields)

	if err != nil {
		return nil, errors.Wrap(err, "unmarshal abi json error")
	}

	contract := &contractImpl{
		name: name,
	}

	for _, field := range fields {
		switch field.Type {
		case abi.JSONTypeConstructor:
			contract.constructor = toExecSite(&field)
		case abi.JSONTypeFunc:
			if *field.StateMutability == abi.StateMutabilityPure || *field.StateMutability == abi.StateMutabilityView {
				contract.calls[field.Name] = toCallSite(&field)
			} else {
				contract.execs[field.Name] = toExecSite(&field)
			}
		case abi.JSONTypeReceive:
			contract.receive = toExecSite(&field)
		case abi.JSONTypeFallback:
			contract.fallback = toExecSite(&field)
		case abi.JSONTypeEvent:
			contract.events[field.Name] = toEventSite(&field)
		case abi.JSONTypeError:
			contract.errors[field.Name] = toErrorSite(&field)
		}
	}

	return contract, nil
}

func (impl *contractImpl) Name() string {
	return impl.name
}

func (impl *contractImpl) Connect(address string, provider client.Provider) {
	impl.address = address
	impl.provider = provider
}

func (impl *contractImpl) Call(ctx context.Context, method string, result interface{}, args ...interface{}) error {
	return nil
}

func (impl *contractImpl) Exec(ctx context.Context, signer signer.Signer, method string, args ...interface{}) error {
	return nil
}
