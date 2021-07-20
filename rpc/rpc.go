package rpc

import (
	"context"

	"github.com/libs4go/fixed"
)

// Block eth block object
type Block struct {
	Number           string         `json:"number"`
	Hash             string         `json:"hash"`
	Parent           string         `json:"parentHash"`
	Nonce            string         `json:"nonce"`
	SHA3Uncles       string         `json:"sha3Uncles"`
	LogsBloom        string         `json:"logsBloom"`
	TransactionsRoot string         `json:"transactionsRoot"`
	StateRoot        string         `json:"stateRoot"`
	ReceiptsRoot     string         `json:"receiptsRoot"`
	Miner            string         `json:"miner"`
	Difficulty       string         `json:"difficulty"`
	TotalDifficulty  string         `json:"totalDifficulty"`
	ExtraData        string         `json:"extraData"`
	Size             string         `json:"size"`
	GasLimit         string         `json:"gasLimit"`
	GasUsed          string         `json:"gasUsed"`
	Timestamp        string         `json:"timestamp"`
	Transactions     []*Transaction `json:"transactions"`
	Uncles           []string       `json:"uncles"`
}

// Transaction .
type Transaction struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	Gas              string `json:"gas"`
	Input            string `json:"input"`
}

// TransactionReceipt .
type TransactionReceipt struct {
	Hash              string        `json:"transactionHash"`
	BlockHash         string        `json:"blockHash"`
	BlockNumber       string        `json:"blockNumber"`
	TransactionIndex  string        `json:"transactionIndex"`
	CumulativeGasUsed string        `json:"cumulativeGasUsed"`
	GasUsed           string        `json:"gasUsed"`
	ContractAddress   string        `json:"contractAddress"`
	Logs              []interface{} `json:"logs"`
	LogsBloom         string        `json:"logsBloom"`
	Status            string        `json:"status"`
}

// CallSite .
type CallSite struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Value    string `json:"value,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Gas      string `json:"gas,omitempty"`
	Data     string `json:"data,omitempty"`
}

// Provider rpc provider
type Provider interface {
	Nonce(ctx context.Context, address string) (uint64, error)
	GetBalance(ctx context.Context, address string) (*fixed.Number, error)
	BestBlockNumber(ctx context.Context) (int64, error)
	Call(ctx context.Context, callsite *CallSite) (val string, err error)
	GetBlockByNumber(ctx context.Context, number int64) (val *Block, err error)
	GetTransactionByHash(ctx context.Context, tx string) (val *Transaction, err error)
	SendRawTransaction(ctx context.Context, tx []byte) (val string, err error)
	GetTransactionReceipt(ctx context.Context, tx string) (val *TransactionReceipt, err error)
	BalanceOfAsset(ctx context.Context, address string, asset string, decimals int) (*fixed.Number, error)
	DecimalsOfAsset(ctx context.Context, asset string) (int, error)
	SuggestGasPrice(ctx context.Context) (*fixed.Number, error)
}
