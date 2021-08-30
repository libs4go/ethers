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

type Foo interface {
	DAO(ctx context.Context) (ret0 address.Address, err error)

	Approve(ctx context.Context, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error)

	BalanceOf(ctx context.Context, owner address.Address) (ret0 *big.Int, err error)

	Burn(ctx context.Context, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error)

	BurnRequire(ctx context.Context, tokenId *big.Int) (ret0 *big.Int, err error)

	CommissionRate(ctx context.Context) (ret0 *big.Int, err error)

	Data(ctx context.Context, tokenId *big.Int) (ret0 *CurveNFT, err error)

	Deposit(ctx context.Context, recipient address.Address, asset address.Address, amount *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error)

	GetApproved(ctx context.Context, tokenId *big.Int) (ret0 address.Address, err error)

	Hello(ctx context.Context, tokenId [20][]*big.Int, nft []*CurveNFT, nfts [][2]*CurveNFT, ops ...abi.Op) (ret0 abi.Transaction, err error)

	IsApprovedForAll(ctx context.Context, owner address.Address, operator address.Address) (ret0 bool, err error)

	Name(ctx context.Context) (ret0 string, err error)

	Owner(ctx context.Context) (ret0 address.Address, err error)

	OwnerOf(ctx context.Context, tokenId *big.Int) (ret0 address.Address, err error)

	RenounceOwnership(ctx context.Context, ops ...abi.Op) (ret0 abi.Transaction, err error)

	SafeTransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error)

	SafeTransferFrom1(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, _data []byte, ops ...abi.Op) (ret0 abi.Transaction, err error)

	SetApprovalForAll(ctx context.Context, operator address.Address, approved bool, ops ...abi.Op) (ret0 abi.Transaction, err error)

	SupportsInterface(ctx context.Context, interfaceId [4]byte) (ret0 bool, err error)

	Symbol(ctx context.Context) (ret0 string, err error)

	TokenByIndex(ctx context.Context, index *big.Int) (ret0 *big.Int, err error)

	TokenOfOwnerByIndex(ctx context.Context, owner address.Address, index *big.Int) (ret0 *big.Int, err error)

	TokenURI(ctx context.Context, tokenId *big.Int) (ret0 string, err error)

	TotalSupply(ctx context.Context) (ret0 *big.Int, err error)

	TransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error)

	TransferOwnership(ctx context.Context, newOwner address.Address, ops ...abi.Op) (ret0 abi.Transaction, err error)

	Usd(ctx context.Context) (ret0 address.Address, err error)

	Withdraw(ctx context.Context, recipient address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error)

	WithdrawAmount(ctx context.Context, tokenId *big.Int) (ret0 *big.Int, err error)

	Withdrawable(ctx context.Context, tokenId *big.Int) (ret0 bool, err error)
}

type implFoo struct {
	contract  abi.Contract
	client    client.Provider
	signer    signer.Signer
	recipient string
}

func (impl *implFoo) DAO(ctx context.Context) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "98fabd3a")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Approve(ctx context.Context, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "095ea7b3")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) BalanceOf(ctx context.Context, owner address.Address) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "70a08231")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Burn(ctx context.Context, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "42966c68")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) BurnRequire(ctx context.Context, tokenId *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "c3a95aeb")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) CommissionRate(ctx context.Context) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "5ea1d6f8")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Data(ctx context.Context, tokenId *big.Int) (ret0 *CurveNFT, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "f0ba8440")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Deposit(ctx context.Context, recipient address.Address, asset address.Address, amount *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "8340f549")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) GetApproved(ctx context.Context, tokenId *big.Int) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "081812fc")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Hello(ctx context.Context, tokenId [20][]*big.Int, nft []*CurveNFT, nfts [][2]*CurveNFT, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "039ef37b")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) IsApprovedForAll(ctx context.Context, owner address.Address, operator address.Address) (ret0 bool, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "e985e9c5")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Name(ctx context.Context) (ret0 string, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "06fdde03")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Owner(ctx context.Context) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "8da5cb5b")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) OwnerOf(ctx context.Context, tokenId *big.Int) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "6352211e")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) RenounceOwnership(ctx context.Context, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "715018a6")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) SafeTransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "42842e0e")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) SafeTransferFrom1(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, _data []byte, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "b88d4fde")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) SetApprovalForAll(ctx context.Context, operator address.Address, approved bool, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "a22cb465")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) SupportsInterface(ctx context.Context, interfaceId [4]byte) (ret0 bool, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "01ffc9a7")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Symbol(ctx context.Context) (ret0 string, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "95d89b41")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) TokenByIndex(ctx context.Context, index *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "4f6ccce7")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) TokenOfOwnerByIndex(ctx context.Context, owner address.Address, index *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "2f745c59")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) TokenURI(ctx context.Context, tokenId *big.Int) (ret0 string, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "c87b56dd")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) TotalSupply(ctx context.Context) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "18160ddd")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) TransferFrom(ctx context.Context, from address.Address, to address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "23b872dd")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) TransferOwnership(ctx context.Context, newOwner address.Address, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "f2fde38b")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) Usd(ctx context.Context) (ret0 address.Address, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "d63a6ccd")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Withdraw(ctx context.Context, recipient address.Address, tokenId *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "f3fef3a3")

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
	callOps, err = abi.MakeCallOps(ctx, impl.client, impl.signer, ops)

	if err != nil {
		return
	}

	ret0, err = abi.MakeTransaction(ctx, impl.client, impl.signer, callOps, impl.recipient, buff)

	return

}

func (impl *implFoo) WithdrawAmount(ctx context.Context, tokenId *big.Int) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "0562b9f7")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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

func (impl *implFoo) Withdrawable(ctx context.Context, tokenId *big.Int) (ret0 bool, err error) {
	f, ok := abi.TryGetFunc(impl.contract, "f11988e0")

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
		To:   impl.recipient,
		Data: hex.EncodeToString(buff),
	}

	var ret string

	ret, err = impl.client.Call(ctx, callSite)

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
