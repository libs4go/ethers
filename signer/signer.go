package signer

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/libs4go/crypto/bip32"
	"github.com/libs4go/crypto/bip39"
	ecdsax "github.com/libs4go/crypto/ecdsa"
	"github.com/libs4go/crypto/ecdsa/recoverable"
	ellipticx "github.com/libs4go/crypto/elliptic"
	"github.com/libs4go/errors"
	"github.com/libs4go/ethers"
	"github.com/libs4go/ethers/address"
	"github.com/libs4go/ethers/eip712"
)

type hdWalletSigner struct {
	addr       string
	privateKey *ecdsa.PrivateKey
}

type keyParam struct {
}

func (*keyParam) Curve() elliptic.Curve {
	return ellipticx.SECP256K1()
}

func NewHDWallet(length int) (string, error) {

	entropy, err := bip39.NewEntropy(length * 8)

	if err != nil {
		return "", errors.Wrap(err, "Create Entropy(%d) error", length)
	}

	return bip39.NewMnemonic(entropy, bip39.ENUS())
}

func OpenHDWallet(mnemonic string, bip44Path string) (ethers.Signer, error) {
	// check mnemonic
	_, err := bip39.MnemonicToByteArray(mnemonic, bip39.ENUS())

	if err != nil {
		return nil, errors.Wrap(err, "invalid mnemonic")
	}

	masterkey, err := bip32.FromMnemonic(&keyParam{}, mnemonic, "")

	if err != nil {
		return nil, errors.Wrap(err, "create master key from mnemonic error")
	}

	privateKeyBytes, err := bip32.DriveFrom(masterkey, bip44Path)

	if err != nil {
		return nil, err
	}

	privateKey := ecdsax.BytesToPrivateKey(privateKeyBytes, ellipticx.SECP256K1())

	return &hdWalletSigner{
		addr:       address.BytesToAddress(ecdsax.PublicKeyBytes(&privateKey.PublicKey)).Hex(),
		privateKey: privateKey,
	}, nil
}

func (wallet *hdWalletSigner) Addresss() string {
	return wallet.addr
}

func (wallet *hdWalletSigner) SignTypedData(typedData *ethers.TypedData) ([]byte, error) {
	return eip712.Sign(wallet.privateKey, (*eip712.TypedData)(typedData))
}

func (wallet *hdWalletSigner) SignTransaction(tx *ethers.Transaction) error {
	r, s, v, err := recoverable.Sign(wallet.privateKey, tx.SignHash(), true)

	if err != nil {
		return err
	}

	tx.R = r
	tx.S = s
	tx.V = v

	return nil
}
