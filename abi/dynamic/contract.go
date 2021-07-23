package dynamic

import (
	"context"

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

// Load JSON abi and set contract name
func Load(name string, abi []byte) (Contract, error) {
	return nil, nil
}
