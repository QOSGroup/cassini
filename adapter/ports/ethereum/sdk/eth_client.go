package sdk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/QOSGroup/cassini/log"
	"github.com/ethereum/go-ethereum/rpc"
)

var ethCLI *EthClient

func init() {
	ethCLI = &EthClient{}
	ethCLI.ctx, ethCLI.cancel =
		context.WithCancel(context.Background())
	// defer cancel()

	// client, err := rpc.DialContext(ctx, "http://192.168.1.178:8545/")
	var err error
	ethCLI.client, err = rpc.DialContext(ethCLI.ctx,
		"https://kovan.infura.io/v3/fb298d4afd444cd5b7c5703b99d51f05")

	if err != nil {
		log.Error("connect to ethereum error: ", err)
	}
}

// EthClient ethereum json-rpc client
type EthClient struct {
	client *rpc.Client
	cancel context.CancelFunc
	ctx    context.Context
}

// JSONRPC wrapped the rpc request
type JSONRPC struct {
	JSONRPC string        `json:"jsonrpc,omitempty"`
	Method  string        `json:"method,omitempty"`
	Params  []interface{} `json:"params,omitempty"`
	ID      int           `json:"id,omitempty"`
}

// Transaction wrap transaction in block
type Transaction struct {
	BlockHash        string `json:"blockHash,omitempty"`
	TransactionIndex string `json:"transactionIndex,omitempty"`
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
}

// Result returns from ethereum
type Result struct {
	Difficulty   string         `json:"difficulty,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
}

// Close client's connection
func Close() {
	ethCLI.cancel()
}

// EthBlockNumber json-rpc for eth_blockNumber
func EthBlockNumber() (height int64, err error) {
	var resp string
	err = ethCLI.client.CallContext(ethCLI.ctx, &resp,
		"eth_blockNumber")
	if err != nil {
		return
	}
	height, err = strconv.ParseInt(resp[2:], 16, 64)
	return
}

// EthGetBlockByNumber json-rpc for eth_getBlockByNumber
func EthGetBlockByNumber(height int64) (result Result, err error) {
	log.Info("height: ", height)
	number := fmt.Sprintf("0x%s", strconv.FormatInt(height, 16))
	log.Info("number: ", number)
	err = ethCLI.client.CallContext(ethCLI.ctx, &result,
		"eth_getBlockByNumber", number, true)
	return
}

// EthProtocolVersion json-rpc for eth_protocolVersion
func EthProtocolVersion() (result string, err error) {
	err = ethCLI.client.CallContext(ethCLI.ctx, &result,
		"eth_protocolVersion")
	return
}
