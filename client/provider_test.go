package client

import (
	"context"
	"math/big"
	"testing"

	"github.com/libs4go/crypto/ecdsa"
	"github.com/libs4go/crypto/elliptic"
	"github.com/libs4go/ethers/address"
	"github.com/libs4go/ethers/signer"
	"github.com/libs4go/fixed"
	"github.com/libs4go/scf4go"
	_ "github.com/libs4go/scf4go/codec/json" //
	"github.com/libs4go/scf4go/reader/memory"
	"github.com/libs4go/slf4go"
	"github.com/stretchr/testify/require"

	_ "github.com/libs4go/slf4go/backend/console" //
)

var configFile = `
{
    "default": {
		"backend": "console"
	},
	"backend": {
		"console": {
			"formatter": {
				"timestamp":"Mon, 02 Jan 2006 15:04:05 -0700",
				"output":"@t @l @s @m"
			}
		}
	}
}
`

func init() {
	config := scf4go.New()

	err := config.Load(memory.New(memory.Data(configFile, "json")))

	if err != nil {
		panic(err)
	}

	err = slf4go.Config(config)

	if err != nil {
		panic(err)
	}
}

func TestRawTransaction(t *testing.T) {

	s, err := signer.OpenHDWallet("orchard mean picnic worry sleep squeeze auto copy hard eager island entry define dune raise spice steel voice prosper mosquito warm ignore book negative", "m/44'/60'/0'/0/0")

	println("wallet", s.Addresss())

	require.NoError(t, err)

	provider, err := HttpProvider("https://data-seed-prebsc-1-s1.binance.org:8545")

	require.NoError(t, err)

	nonce, err := provider.Nonce(context.Background(), s.Addresss())

	require.NoError(t, err)

	gasPrice, err := fixed.New(18, fixed.Float(0.000000018))

	require.NoError(t, err)

	amount, err := fixed.New(18, fixed.Float(0.1))

	require.NoError(t, err)

	recipient := [20]byte(address.HexToAddress("0x44A347Cf7278685320a05Cb39e903C42e472e262"))

	println("value:", gasPrice.RawValue.String(), amount.RawValue.String())

	tx := &signer.Transaction{
		AccountNonce: nonce,
		Price:        gasPrice.RawValue,
		GasLimit:     big.NewInt(21000),
		Recipient:    &recipient,
		Amount:       amount.RawValue,
	}

	err = s.SignTransaction(tx)

	require.NoError(t, err)

	pk, _, err := ecdsa.Recover(elliptic.SECP256K1(), tx.R, tx.S, tx.V, tx.SignHash())

	require.NoError(t, err)

	println("recover pk", address.FromPublicKey(pk).Hex())

	rawTx, err := tx.Encode()

	require.NoError(t, err)

	txID, err := provider.SendRawTransaction(context.Background(), rawTx)

	require.NoError(t, err)

	println(txID)
}

func TestFetchTransaction(t *testing.T) {
	provider, err := HttpProvider("https://bsc-dataseed1.binance.org/")

	require.NoError(t, err)

	_, err = provider.GetTransactionReceipt(context.Background(), "0x164acf4eb987a486d131d36c36d07c46f4306ee96c29bad9172fd3fa3e7edbcf")

	require.NoError(t, err)

}
