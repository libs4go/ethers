package signer

import (
	"encoding/hex"
	"math/big"

	"github.com/libs4go/encoding/rlp"
	"github.com/libs4go/ethers/eip712"
	"golang.org/x/crypto/sha3"
)

// Transaction ether transaction object
type Transaction struct {
	AccountNonce uint64    `json:"nonce"    gencodec:"required"`
	Price        *big.Int  `json:"gasPrice" gencodec:"required"`
	GasLimit     *big.Int  `json:"gas"      gencodec:"required"`
	Recipient    *[20]byte `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int  `json:"value"    gencodec:"required"`
	Payload      []byte    `json:"input"    gencodec:"required"`
	V            *big.Int  `json:"v" gencodec:"required"`
	R            *big.Int  `json:"r" gencodec:"required"`
	S            *big.Int  `json:"s" gencodec:"required"`
}

// Hash get tx hash string
func (tx *Transaction) Hash() string {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, tx)
	return "0x" + hex.EncodeToString(hw.Sum(nil))
}

func (tx *Transaction) SignHash() []byte {

	hw := sha3.NewLegacyKeccak256()

	rlp.Encode(hw, []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
	})

	return hw.Sum(nil)
}

// Encode encode tx to raw transaction bytes
func (tx *Transaction) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

// TypedData ...
type TypedData eip712.TypedData

// VerifyTypedData verify eip712 sig, and return signer address ...
func VerifyTypedData(typedData *TypedData, sig []byte) (string, error) {
	return eip712.Recover((*eip712.TypedData)(typedData), sig)
}

// Signer the ethers signer ....
type Signer interface {
	// The signer ether address
	Addresss() string
	// SignTypedData implement eip712 sign ...
	SignTypedData(typedData *TypedData) ([]byte, error)
	// SignTransaction sign ether transaction
	SignTransaction(tx *Transaction) error
}
