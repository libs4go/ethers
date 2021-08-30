package testdata

import (
	"context"
	"encoding/hex"

	"github.com/libs4go/errors"
	"github.com/libs4go/ethers/abi"
	"github.com/libs4go/ethers/abi/binding"
	"github.com/libs4go/ethers/client"
	"github.com/libs4go/ethers/signer"

	"math/big"

	"github.com/libs4go/ethers/address"
)

// Generated tuple "CurveNFT" stub code , do not modify manually
type CurveNFT struct {
	CommissionAmount *big.Int

	Created *big.Int

	Deposit address.Address

	DepositAmount *big.Int

	Id *big.Int
}

// type CurveUSDVault interface {
//
// 	DAO(ctx context.Context, )(ret0 address.Address, err error)
//
// 	Approve(ctx context.Context, to address.Address, tokenId *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	BalanceOf(ctx context.Context, owner address.Address)(ret0 *big.Int, err error)
//
// 	Burn(ctx context.Context, tokenId *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	BurnRequire(ctx context.Context, tokenId *big.Int)(ret0 *big.Int, err error)
//
// 	CommissionRate(ctx context.Context, )(ret0 *big.Int, err error)
//
// 	Data(ctx context.Context, tokenId *big.Int)(ret0 *CurveNFT, err error)
//
// 	Deposit(ctx context.Context, recipient address.Address, asset address.Address, amount *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	GetApproved(ctx context.Context, tokenId *big.Int)(ret0 address.Address, err error)
//
// 	Hello(ctx context.Context, tokenId [20][]*big.Int, nft []*CurveNFT, nfts [][2]*CurveNFT, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	IsApprovedForAll(ctx context.Context, owner address.Address, operator address.Address)(ret0 bool, err error)
//
// 	Name(ctx context.Context, )(ret0 string, err error)
//
// 	Owner(ctx context.Context, )(ret0 address.Address, err error)
//
// 	OwnerOf(ctx context.Context, tokenId *big.Int)(ret0 address.Address, err error)
//
// 	RenounceOwnership(ctx context.Context, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	SafeTransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	SafeTransferFrom1(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, _data []byte, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	SetApprovalForAll(ctx context.Context, operator address.Address, approved bool, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	SupportsInterface(ctx context.Context, interfaceId [4]byte)(ret0 bool, err error)
//
// 	Symbol(ctx context.Context, )(ret0 string, err error)
//
// 	TokenByIndex(ctx context.Context, index *big.Int)(ret0 *big.Int, err error)
//
// 	TokenOfOwnerByIndex(ctx context.Context, owner address.Address, index *big.Int)(ret0 *big.Int, err error)
//
// 	TokenURI(ctx context.Context, tokenId *big.Int)(ret0 string, err error)
//
// 	TotalSupply(ctx context.Context, )(ret0 *big.Int, err error)
//
// 	TransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	TransferOwnership(ctx context.Context, newOwner address.Address, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	Usd(ctx context.Context, )(ret0 address.Address, err error)
//
// 	Withdraw(ctx context.Context, recipient address.Address, tokenId *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	WithdrawAmount(ctx context.Context, tokenId *big.Int)(ret0 *big.Int, err error)
//
// 	Withdrawable(ctx context.Context, tokenId *big.Int)(ret0 bool, err error)
//
// }

type CurveUSDVault struct {
	Contract  abi.Contract
	Client    client.Provider
	Signer    signer.Signer
	Recipient string
}

func (impl *CurveUSDVault) DAO(ctx context.Context) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "98fabd3a")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func DAO not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Approve(ctx context.Context, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "095ea7b3")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Approve not found")
		return
	}

	var buff []byte

	buff, err = f.Call(to, tokenId)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) BalanceOf(ctx context.Context, owner address.Address) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "70a08231")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func BalanceOf not found")
		return
	}

	var buff []byte

	buff, err = f.Call(owner)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Burn(ctx context.Context, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "42966c68")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Burn not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) BurnRequire(ctx context.Context, tokenId *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "c3a95aeb")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func BurnRequire not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) CommissionRate(ctx context.Context) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "5ea1d6f8")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func CommissionRate not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Data(ctx context.Context, tokenId *big.Int) (ret0 *CurveNFT, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "f0ba8440")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Data not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Deposit(ctx context.Context, recipient address.Address, asset address.Address, amount *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "8340f549")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Deposit not found")
		return
	}

	var buff []byte

	buff, err = f.Call(recipient, asset, amount)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) GetApproved(ctx context.Context, tokenId *big.Int) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "081812fc")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func GetApproved not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Hello(ctx context.Context, tokenId [20][]*big.Int, nft []*CurveNFT, nfts [][2]*CurveNFT, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "039ef37b")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Hello not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId, nft, nfts)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) IsApprovedForAll(ctx context.Context, owner address.Address, operator address.Address) (ret0 bool, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "e985e9c5")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func IsApprovedForAll not found")
		return
	}

	var buff []byte

	buff, err = f.Call(owner, operator)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Name(ctx context.Context) (ret0 string, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "06fdde03")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Name not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Owner(ctx context.Context) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "8da5cb5b")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Owner not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) OwnerOf(ctx context.Context, tokenId *big.Int) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "6352211e")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func OwnerOf not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) RenounceOwnership(ctx context.Context, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "715018a6")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func RenounceOwnership not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) SafeTransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "42842e0e")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func SafeTransferFrom not found")
		return
	}

	var buff []byte

	buff, err = f.Call(from, to, tokenId)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) SafeTransferFrom1(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, _data []byte, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "b88d4fde")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func SafeTransferFrom1 not found")
		return
	}

	var buff []byte

	buff, err = f.Call(from, to, tokenId, _data)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) SetApprovalForAll(ctx context.Context, operator address.Address, approved bool, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "a22cb465")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func SetApprovalForAll not found")
		return
	}

	var buff []byte

	buff, err = f.Call(operator, approved)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) SupportsInterface(ctx context.Context, interfaceId [4]byte) (ret0 bool, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "01ffc9a7")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func SupportsInterface not found")
		return
	}

	var buff []byte

	buff, err = f.Call(interfaceId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Symbol(ctx context.Context) (ret0 string, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "95d89b41")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Symbol not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) TokenByIndex(ctx context.Context, index *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "4f6ccce7")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TokenByIndex not found")
		return
	}

	var buff []byte

	buff, err = f.Call(index)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) TokenOfOwnerByIndex(ctx context.Context, owner address.Address, index *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "2f745c59")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TokenOfOwnerByIndex not found")
		return
	}

	var buff []byte

	buff, err = f.Call(owner, index)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) TokenURI(ctx context.Context, tokenId *big.Int) (ret0 string, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "c87b56dd")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TokenURI not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) TotalSupply(ctx context.Context) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "18160ddd")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TotalSupply not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) TransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "23b872dd")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TransferFrom not found")
		return
	}

	var buff []byte

	buff, err = f.Call(from, to, tokenId)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) TransferOwnership(ctx context.Context, newOwner address.Address, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "f2fde38b")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TransferOwnership not found")
		return
	}

	var buff []byte

	buff, err = f.Call(newOwner)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) Usd(ctx context.Context) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "d63a6ccd")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Usd not found")
		return
	}

	var buff []byte

	buff, err = f.Call()

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Withdraw(ctx context.Context, recipient address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "f3fef3a3")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Withdraw not found")
		return
	}

	var buff []byte

	buff, err = f.Call(recipient, tokenId)

	if err != nil {
		return
	}

	var callOps *abi.CallOps
	callOps, err = abi.MakeCallOps(ctx, impl.Client, impl.Signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.Client, impl.Signer, callOps, impl.Recipient, buff)

	return

}

func (impl *CurveUSDVault) WithdrawAmount(ctx context.Context, tokenId *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "0562b9f7")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func WithdrawAmount not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}

func (impl *CurveUSDVault) Withdrawable(ctx context.Context, tokenId *big.Int) (ret0 bool, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "f11988e0")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Withdrawable not found")
		return
	}

	var buff []byte

	buff, err = f.Call(tokenId)

	if err != nil {
		return
	}

	callSite := &client.CallSite{
		To:   impl.Recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.Client.Call(ctx, callSite)

	if err != nil {
		return
	}

	buff, err = hex.DecodeString(ret)

	if err != nil {
		return
	}

	_, err = f.Return(buff, []interface{}{&ret0})

	return

}
