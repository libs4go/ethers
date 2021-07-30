package abi

import "github.com/libs4go/errors"

// ScopeOfAPIError .
const errVendor = "ethers-abi"

// errors
var (
	ErrBits       = errors.New("integer type bits out of range or is not a multiple of 8", errors.WithVendor(errVendor), errors.WithCode(-1))
	ErrValue      = errors.New("encode value type error", errors.WithVendor(errVendor), errors.WithCode(-2))
	ErrFixedBytes = errors.New("fixed bytes length mismatch", errors.WithVendor(errVendor), errors.WithCode(-3))
	ErrLength     = errors.New("length error", errors.WithVendor(errVendor), errors.WithCode(-4))
)
