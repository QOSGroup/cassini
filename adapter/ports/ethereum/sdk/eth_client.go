package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/QOSGroup/cassini/log"
)

var cli *EthClient

func init() {
	cli = &EthClient{
		client:  &http.Client{},
		url:     "https://kovan.infura.io/v3/fb298d4afd444cd5b7c5703b99d51f05",
		jsonrpc: "2.0",
		mux:     new(sync.Mutex)}
	cli.ctx, cli.cancel = context.WithCancel(context.Background())
}

// JSONRPC rpc request
type JSONRPC struct {
	JSONRPC string        `json:"jsonrpc,omitempty"`
	Method  string        `json:"method,omitempty"`
	Params  []interface{} `json:"params,omitempty"`
	ID      uint8         `json:"id,omitempty"`
}

// Response rpc response
type Response struct {
	JSONRPCVersion string       `json:"jsonrpc"`
	ID             uint8        `json:"id"`
	Result         string       `json:"result"`
	Error          *ResultError `json:"error"`
}

// ResultError error message
type ResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Transaction block tx
type Transaction struct {
	BlockHash        string `json:"blockHash,omitempty"`
	TransactionIndex string `json:"transactionIndex,omitempty"`
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
	Value            string `json:"value,omitempty"`

	BlockNumber string `json:"blockNumber,omitempty"`
	ChainID     string `json:"chainId,omitempty"`
	//     Condition
	//     Creates
	Gas       string `json:"gas,omitempty"`
	GasPrice  string `json:"gasPrice,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Input     string `json:"input,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	R         string `json:"r,omitempty"`
	Raw       string `json:"raw,omitempty"`
	S         string `json:"s,omitempty"`
	StandardV string `standardV:"value,omitempty"`
	V         string `json:"v,omitempty"`
}

// Block returns from ethereum
type Block struct {
	Number       string         `json:"number,omitempty"`
	Difficulty   string         `json:"difficulty,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
}

// ResponseBlock returns block from ethereum
type ResponseBlock struct {
	JSONRPCVersion string `json:"jsonrpc"`
	ID             uint8  `json:"id"`
	Result         *Block `json:"result"`
}

// Close client's connection
func Close() {
	cli.cancel()
}

// EthBlockNumberResponse json-rpc response string for eth_blockNumber
func EthBlockNumberResponse() (response string, err error) {
	response, err = cli.call("eth_blockNumber")
	log.Info("response: ", response)
	return
}

// EthBlockNumber return block number in hex string
func EthBlockNumber() (height string, err error) {
	resp, err := EthBlockNumberResponse()
	if err != nil {
		return
	}
	var response Response
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return
	}
	height = response.Result
	log.Debug("height: ", height)
	if response.Error != nil {
		err = fmt.Errorf("code: %d; message: %s",
			response.Error.Code, response.Error.Message)
		return
	}
	return
}

// EthBlockNumberInt64 return block number in int64
func EthBlockNumberInt64() (height int64, err error) {
	h, err := EthBlockNumber()
	if err != nil {
		return
	}
	height, err = strconv.ParseInt(h[2:], 16, 64)
	return
}

// EthGetBlockByNumberResponse json-rpc for eth_getBlockByNumber
func EthGetBlockByNumberResponse(height int64) (response string, err error) {
	number := fmt.Sprintf("0x%s", strconv.FormatInt(height, 16))
	log.Infof("height: %d number: %s", height, number)
	response, err = cli.call("eth_getBlockByNumber", number, true)
	log.Info("response: ", response)
	return
}

// EthGetBlockByNumber json-rpc for eth_getBlockByNumber
func EthGetBlockByNumber(height int64) (*Block, error) {
	resp, err := EthGetBlockByNumberResponse(height)
	if err != nil {
		return nil, err
	}
	var response ResponseBlock
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// EthProtocolVersionResponse json-rpc for eth_protocolVersion
func EthProtocolVersionResponse() (response string, err error) {
	response, err = cli.call("eth_protocolVersion")
	return
}

// EthClient ethereum json-rpc client
type EthClient struct {
	client  *http.Client
	cancel  context.CancelFunc
	ctx     context.Context
	url     string
	jsonrpc string
	id      uint8
	mux     *sync.Mutex
}

func (c *EthClient) call(method string,
	params ...interface{}) (response string, err error) {
	rpc := c.newRPC(method, params...)
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(rpc)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, c.url, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response = string(body)
	return
}

func (c *EthClient) increaseID() uint8 {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.id++
	if c.id == 0 {
		c.id++
	}
	id := c.id
	return id
}

func (c *EthClient) newRPC(method string, params ...interface{}) (rpc *JSONRPC) {
	id := c.increaseID()
	rpc = &JSONRPC{
		JSONRPC: c.jsonrpc,
		ID:      id,
		Method:  method}
	if len(params) > 0 {
		log.Debugf("call %s len(params): %d", method, len(params))
		rpc.Params = append(rpc.Params, params...)
	} else {
		log.Debugf("call %s params is 0", method)
	}
	return
}
