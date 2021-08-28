package binding

import "github.com/libs4go/ethers/abi"

// ABI binder context
type Binder interface {
	RegisterTuple(tuple *Tuple)
	GetTuple(name string) (*Tuple, bool)
	BeginContract(name string, abi []byte)
	Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder)
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

func (impl *symbolsImpl) Func(name string, selector string, inputs []abi.Encoder, outputs []abi.Encoder) {

}

func (impl *symbolsImpl) EndContract() {

}

func NewGen() Binder {
	return nil
}
