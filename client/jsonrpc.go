package client

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/libs4go/errors"
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
func (client *jsonrpcProvider) BlockNumber(ctx context.Context) (uint64, error) {

	var data string

	err := client.rpcCall(ctx, "eth_blockNumber", &data)

	if err != nil {
		return 0, err
	}

	val, err := fixed.New(0, fixed.HexRawValue(data))

	if err != nil {
		return 0, errors.Wrap(err, "decode %s error", data)
	}

	return uint64(val.RawValue.Int64()), nil
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

func (client *jsonrpcProvider) GetBlockTransactionCountByHash(ctx context.Context, blockHash string) (uint64, error) {
	var data string

	err := client.rpcCall(ctx, "eth_getBlockTransactionCountByHash", &data, blockHash)

	if err != nil {
		return 0, err
	}

	val, err := fixed.New(0, fixed.HexRawValue(data))

	if err != nil {
		return 0, errors.Wrap(err, "decode %s error", data)
	}

	return uint64(val.RawValue.Int64()), nil
}

func (client *jsonrpcProvider) GetBlockTransactionCountByNumber(ctx context.Context, number uint64) (uint64, error) {
	var data string

	err := client.rpcCall(ctx, "eth_getBlockTransactionCountByHash", &data, fmt.Sprintf("0x%x", number))

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
func (client *jsonrpcProvider) GetBlockByNumber(ctx context.Context, number uint64, full bool) (val *Block, err error) {

	err = client.rpcCall(ctx, "eth_getBlockByNumber", &val, fmt.Sprintf("0x%x", number), full)

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

// SuggestGasPrice .
func (client *jsonrpcProvider) GasPrice(ctx context.Context) (*fixed.Number, error) {
	var val string

	err := client.rpcCall(ctx, "eth_gasPrice", &val)
	if err != nil {
		return nil, err
	}

	return fixed.New(18, fixed.HexRawValue(val))
}

func (client *jsonrpcProvider) GetBlockByHash(ctx context.Context, blockHash string, full bool) (val *Block, err error) {
	err = client.rpcCall(ctx, "eth_getBlockByNumber", &val, blockHash, full)

	return
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
