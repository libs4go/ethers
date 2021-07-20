package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/libs4go/ethers/address"
	"github.com/libs4go/fixed"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
	s, err := OpenHDWallet("orchard mean picnic worry sleep squeeze auto copy hard eager island entry define dune raise spice steel voice prosper mosquito warm ignore book negative", "m/44'/60'/0'/0/0")

	require.NoError(t, err)

	gasPrice, err := fixed.New(18, fixed.Float(0.000000018))

	require.NoError(t, err)

	amount, err := fixed.New(18, fixed.Float(0.01))

	require.NoError(t, err)

	println(gasPrice.String())

	println(fmt.Sprintf("%.10f", gasPrice.Float()))

	recipient := [20]byte(address.HexToAddress("0xaa25aa7a19f9c426e07dee59b12f944f4d9f1dd3"))

	tx := &Transaction{
		AccountNonce: 0,
		Price:        gasPrice.RawValue,
		GasLimit:     big.NewInt(21000),
		Recipient:    &recipient,
		Amount:       amount.RawValue,
	}

	err = s.SignTransaction(tx)

	require.NoError(t, err)

	println(tx.Hash())

	data, err := json.Marshal(tx)

	require.NoError(t, err)

	println(string(data))

	rawTx, err := tx.Encode()

	require.NoError(t, err)

	println(hex.EncodeToString(rawTx))
}
