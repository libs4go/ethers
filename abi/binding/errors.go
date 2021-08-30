package binding

import "github.com/libs4go/errors"

// ScopeOfAPIError .
const errVendor = "ethers-abi-binding"

// errors
var (
	ErrBinding = errors.New("Binding internal error", errors.WithVendor(errVendor), errors.WithCode(-1))
)
