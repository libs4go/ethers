package abi

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"

	"github.com/libs4go/ethers/address"
	"github.com/libs4go/ethers/client"
	"github.com/libs4go/ethers/signer"
	"github.com/libs4go/fixed"
	"github.com/libs4go/slf4go"
)

type Func interface {
	Selector() []byte
	// Call generate call bytes
	Call(params ...interface{}) ([]byte, error)
	// Return unmarshal return bytes
	Return(data []byte, values interface{}) (uint, error)
}

type Contract interface {
	Select(selector string) (Func, bool)
}

func TryGetFunc(contract Contract, signature string) (Func, bool) {
	return contract.Select(hex.EncodeToString(Selector(signature)))
}

type CallOps struct {
	GasLimit *big.Int
	GasPrice *big.Int
	Nonce    *big.Int
	Amount   *big.Int
}

type Op func(ops *CallOps)

func WithGasPrice(value *fixed.Number) Op {
	return func(ops *CallOps) {
		ops.GasPrice = value.RawValue
	}
}

func WithGasLimits(value *big.Int) Op {
	return func(ops *CallOps) {
		ops.GasLimit = value
	}
}

func WithNonce(value uint64) Op {
	return func(ops *CallOps) {
		ops.Nonce = new(big.Int).SetUint64(value)
	}
}

func WithAmount(value *fixed.Number) Op {
	return func(ops *CallOps) {
		ops.Amount = value.RawValue
	}
}

func MakeCallOps(ctx context.Context, client client.Provider, signer signer.Signer, ops []Op) (*CallOps, error) {
	callOps := &CallOps{}

	for _, op := range ops {
		op(callOps)
	}

	if callOps.GasLimit == nil {
		WithGasLimits(big.NewInt(21000))(callOps)
	}

	if callOps.GasPrice == nil {
		gasPrice, err := client.GasPrice(ctx)

		if err != nil {
			return nil, err
		}

		WithGasPrice(gasPrice)(callOps)
	}

	if callOps.Nonce == nil {
		if signer != nil {
			nonce, err := client.Nonce(ctx, signer.Addresss())

			if err != nil {
				return nil, err
			}

			WithNonce(nonce)(callOps)
		} else {
			WithNonce(0)(callOps)
		}
	}

	if callOps.Amount == nil {
		WithAmount(&fixed.Number{RawValue: big.NewInt(0), Decimals: 18})(callOps)
	}

	return callOps, nil
}

// Transaction
type Transaction interface {
	Close()
	TX() string
	Receipt() <-chan *TransactionReceipt
}

type transactionImpl struct {
	slf4go.Logger
	txID        string
	receiptChan chan *TransactionReceipt
	client      client.Provider
	once        sync.Once
	cancel      context.CancelFunc
	ctx         context.Context
}

func newTransaction(ctx context.Context, txID string, client client.Provider) Transaction {

	newCTX, cancel := context.WithCancel(ctx)

	return &transactionImpl{
		Logger:      slf4go.Get("ethers-abi-contract"),
		txID:        txID,
		receiptChan: make(chan *TransactionReceipt),
		client:      client,
		ctx:         newCTX,
		cancel:      cancel,
	}
}

func (impl *transactionImpl) TX() string {
	return impl.txID
}

func (impl *transactionImpl) Close() {
	close(impl.receiptChan)
	impl.cancel()
}

func (impl *transactionImpl) doReceipt() {
	defer func() {
		if e := recover(); e != nil {
			impl.E("receipt err {@e}", e)
		}
	}()

	receipt, err := impl.client.GetTransactionReceipt(impl.ctx, impl.txID)

	impl.receiptChan <- &TransactionReceipt{
		Error: err,
		Data:  receipt,
	}
}

func (impl *transactionImpl) Receipt() <-chan *TransactionReceipt {
	impl.once.Do(impl.doReceipt)
	return impl.receiptChan
}

func MakeTransaction(ctx context.Context, client client.Provider, s signer.Signer, callOpts *CallOps, recipient string, data []byte) (Transaction, error) {

	recipientBytes := [20]byte(address.HexToAddress(recipient))

	tx := &signer.Transaction{
		AccountNonce: 0,
		Price:        callOpts.GasPrice,
		GasLimit:     callOpts.GasLimit,
		Recipient:    &recipientBytes,
		Amount:       callOpts.Amount,
		Payload:      data,
	}

	err := s.SignTransaction(tx)

	if err != nil {
		return nil, err
	}

	rawTx, err := tx.Encode()

	if err != nil {
		return nil, err
	}

	txID, err := client.SendRawTransaction(ctx, rawTx)

	if err != nil {
		return nil, err
	}

	return newTransaction(ctx, txID, client), nil
}

type TransactionReceipt struct {
	Error error
	Data  *client.TransactionReceipt
}
