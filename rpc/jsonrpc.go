package rpc

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/libs4go/errors"
	"github.com/libs4go/ethers/abi"
	"github.com/libs4go/fixed"
	"github.com/libs4go/jsonrpc"
	"github.com/libs4go/jsonrpc/client"
)

type jsonrpcProvider struct {
	client jsonrpc.Client
}

func NewJSONRPCProvider(client jsonrpc.Client) (Provider, error) {
	return &jsonrpcProvider{
		client: client,
	}, nil
}

func (client *jsonrpcProvider) rpcCall(ctx context.Context, method string, result interface{}, args ...interface{}) error {
	return client.client.Call(ctx, method, args...).Join(result)
}

func (client *jsonrpcProvider) GetBalance(ctx context.Context, address string) (*fixed.Number, error) {

	var data string

	err := client.rpcCall(ctx, "eth_getBalance", &data, address, "latest")

	if err != nil {
		return nil, err
	}

	return fixed.New(18, fixed.HexRawValue(data))
}

// BlockNumber get geth last block number
func (client *jsonrpcProvider) BestBlockNumber(ctx context.Context) (int64, error) {

	var data string

	err := client.rpcCall(ctx, "eth_blockNumber", &data)

	if err != nil {
		return 0, err
	}

	val, err := fixed.New(0, fixed.HexRawValue(data))

	if err != nil {
		return 0, errors.Wrap(err, "decode %s error", data)
	}

	return val.RawValue.Int64(), nil
}

// Nonce get address send transactions
func (client *jsonrpcProvider) Nonce(ctx context.Context, address string) (uint64, error) {
	var data string

	err := client.rpcCall(ctx, "eth_getTransactionCount", &data, address, "latest")

	if err != nil {
		return 0, err
	}

	val, err := fixed.New(0, fixed.HexRawValue(data))

	if err != nil {
		return 0, errors.Wrap(err, "decode %s error", data)
	}

	return uint64(val.RawValue.Int64()), nil
}

func (client *jsonrpcProvider) Call(ctx context.Context, callsite *CallSite) (val string, err error) {

	err = client.rpcCall(ctx, "eth_call", &val, callsite, "latest")

	return
}

// BlockByNumber get block by number
func (client *jsonrpcProvider) GetBlockByNumber(ctx context.Context, number int64) (val *Block, err error) {

	err = client.rpcCall(ctx, "eth_getBlockByNumber", &val, fmt.Sprintf("0x%x", number), true)

	return
}

// GetTransactionByHash get geth last block number
func (client *jsonrpcProvider) GetTransactionByHash(ctx context.Context, tx string) (val *Transaction, err error) {

	err = client.rpcCall(ctx, "eth_getTransactionByHash", &val, tx)

	return
}

// SendRawTransaction .
func (client *jsonrpcProvider) SendRawTransaction(ctx context.Context, tx []byte) (val string, err error) {

	err = client.rpcCall(ctx, "eth_sendRawTransaction", &val, "0x"+hex.EncodeToString(tx))

	return
}

// GetTransactionReceipt ...
func (client *jsonrpcProvider) GetTransactionReceipt(ctx context.Context, tx string) (val *TransactionReceipt, err error) {

	err = client.rpcCall(ctx, "eth_getTransactionReceipt", &val, tx)

	return
}

// BalanceOfAsset .
func (client *jsonrpcProvider) BalanceOfAsset(ctx context.Context, address string, asset string, decimals int) (*fixed.Number, error) {
	data := abi.BalanceOf(address)

	valstr, err := client.Call(ctx, &CallSite{
		To:   asset,
		Data: data,
	})

	if err != nil {
		return nil, err
	}

	return fixed.New(decimals, fixed.HexRawValue(valstr))
}

// GetTokenDecimals .
func (client *jsonrpcProvider) DecimalsOfAsset(ctx context.Context, asset string) (int, error) {
	data := abi.GetDecimals()

	valstr, err := client.Call(ctx, &CallSite{
		To:   asset,
		Data: data,
	})

	if err != nil {
		return 0, err
	}

	val, err := fixed.New(0, fixed.HexRawValue(valstr))

	if err != nil {
		return 0, errors.Wrap(err, "decode hex %s error", valstr)
	}

	return int(val.RawValue.Int64()), nil
}

// SuggestGasPrice .
func (client *jsonrpcProvider) SuggestGasPrice(ctx context.Context) (*fixed.Number, error) {
	var val string

	err := client.rpcCall(ctx, "eth_gasPrice", &val)
	if err != nil {
		return nil, err
	}

	return fixed.New(18, fixed.HexRawValue(val))
}

// HttpProvider create http jsonrpc provider
func HttpProvider(remote string, ops ...client.ClientOpt) (Provider, error) {
	c, err := client.HTTPConnect(remote, ops...)

	if err != nil {
		return nil, err
	}

	return NewJSONRPCProvider(c)
}

func WebsocketProvider(remote string, ops ...client.ClientOpt) (Provider, error) {
	c, err := client.WebSocketConnect(remote, ops...)

	if err != nil {
		return nil, err
	}

	return NewJSONRPCProvider(c)
}
