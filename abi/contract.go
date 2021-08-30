package abi

import "github.com/libs4go/ethers/client"

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

// Transaction
type Transaction interface {
	Close()
	TX() string
	Receipt() <-chan *TransactionReceipt
}

type TransactionReceipt struct {
	Error error
	Data  *client.TransactionReceipt
}
