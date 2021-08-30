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

// type ERC20 interface {
//
// 	Allowance(ctx context.Context, owner address.Address, spender address.Address)(ret0 *big.Int, err error)
//
// 	Approve(ctx context.Context, spender address.Address, amount *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	BalanceOf(ctx context.Context, account address.Address)(ret0 *big.Int, err error)
//
// 	TotalSupply(ctx context.Context, )(ret0 *big.Int, err error)
//
// 	Transfer(ctx context.Context, recipient address.Address, amount *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// 	TransferFrom(ctx context.Context, sender address.Address, recipient address.Address, amount *big.Int, ops ...abi.Op)(ret0 abi.Transaction, err error)
//
// }

type ERC20 struct {
	Contract  abi.Contract
	Client    client.Provider
	Signer    signer.Signer
	Recipient string
}

func (impl *ERC20) Allowance(ctx context.Context, owner address.Address, spender address.Address) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "dd62ed3e")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Allowance not found")
		return
	}

	var buff []byte

	buff, err = f.Call(owner, spender)

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

func (impl *ERC20) Approve(ctx context.Context, spender address.Address, amount *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "095ea7b3")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Approve not found")
		return
	}

	var buff []byte

	buff, err = f.Call(spender, amount)

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

func (impl *ERC20) BalanceOf(ctx context.Context, account address.Address) (ret0 *big.Int, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "70a08231")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func BalanceOf not found")
		return
	}

	var buff []byte

	buff, err = f.Call(account)

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

func (impl *ERC20) TotalSupply(ctx context.Context) (ret0 *big.Int, err error) {
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

func (impl *ERC20) Transfer(ctx context.Context, recipient address.Address, amount *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "a9059cbb")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func Transfer not found")
		return
	}

	var buff []byte

	buff, err = f.Call(recipient, amount)

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

func (impl *ERC20) TransferFrom(ctx context.Context, sender address.Address, recipient address.Address, amount *big.Int, ops ...abi.Op) (ret0 abi.Transaction, err error) {
	f, ok := abi.TryGetFunc(impl.Contract, "23b872dd")

	if !ok {
		err = errors.Wrap(binding.ErrBinding, "func TransferFrom not found")
		return
	}

	var buff []byte

	buff, err = f.Call(sender, recipient, amount)

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
