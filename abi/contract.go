package abi

type Func interface {
	Selector() []byte
	// Call generate call bytes
	Call(params ...interface{}) ([]byte, error)
	// Return unmarshal return bytes
	Return(data []byte, values interface{}) (uint, error)
}

type Contract interface {
	Func(signature string) (Func, bool)
}
