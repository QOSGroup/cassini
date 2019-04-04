// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
)

// JSONRPC rpc request
type JSONRPC struct {
	JSONRPC string        `json:"jsonrpc,omitempty"`
	Method  string        `json:"method,omitempty"`
	Params  []interface{} `json:"params,omitempty"`
	ID      uint8         `json:"id,omitempty"`
}

// Response rpc response
type Response struct {
	JSONRPCVersion string `json:"jsonrpc"`
	ID             uint8  `json:"id"`
	Result         string `json:"result"`
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

// Result returns from ethereum
type Result struct {
	Difficulty   string         `json:"difficulty,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
}

// BlockResponse returns from ethereum
type BlockResponse struct {
	JSONRPCVersion string  `json:"jsonrpc"`
	ID             uint8   `json:"id"`
	Result         *Result `json:"result"`
}

var id uint8
var txhash = "0x3d947eB8c366D2416468675cEDd00fd311D70dFB"

func main() {

	// s := "0xde0b6b3a7640000"
	// sh, _ := strconv.ParseInt(s[2:], 16, 64)
	// fmt.Println("sh", sh)

	client := &http.Client{}
	var height string
	height = "0xa2fce1"
	// height = "0xa2fd92"
	h, _ := strconv.ParseInt(height[2:], 16, 64)

	go func() {
		for true {
			// resp, err := ethBlockNumber(client)
			// if err != nil {
			// 	fmt.Println("error: ", err)
			// }
			// if !strings.EqualFold(height, resp.Result) {
			// height = resp.Result
			height := fmt.Sprintf("0x%s", strconv.FormatInt(h, 16))
			fmt.Println("======")
			fmt.Println("block number: ", height)
			respBlock, err := ethGetBlockByNumber(client, height)
			if err != nil {
				fmt.Println("error: ", err)
			}
			filterBlock(respBlock)
			// }
			time.Sleep(1 * time.Second)
			if respBlock.Result != nil {
				h++
			}
		}
	}()
	common.KeepRunning(func(sig os.Signal) {
		log.Info("cancel done.")
	})
}

func filterBlock(resp *BlockResponse) {
	if len(resp.Result.Transactions) == 0 {
		return
	}
	for _, tx := range resp.Result.Transactions {
		if strings.EqualFold(tx.To, txhash) {
			fmt.Println("!!!!!!")
			value, _ := strconv.ParseInt(tx.Value[2:], 16, 64)
			fmt.Printf("tx value: %d; %s; %s\n",
				value, tx.Value, tx.BlockNumber)
		}

	}
}

func ethGetBlockByNumber(client *http.Client,
	height string) (*BlockResponse, error) {
	var ps []interface{}
	ps = append(ps, height, true)
	req := JSONRPC{
		Method: "eth_getBlockByNumber",
		Params: ps}
	var resp BlockResponse
	err := call(client, &req, &resp)
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("response jsonrpc: %s; id: %d\n",
			resp.JSONRPCVersion, resp.ID)
	}
	return &resp, err
}

func ethBlockNumber(client *http.Client) (*Response, error) {
	req := JSONRPC{
		Method: "eth_blockNumber"}
	var resp Response
	err := call(client, &req, &resp)
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("response jsonrpc: %s; id: %d; result: %s\n",
			resp.JSONRPCVersion, resp.ID, resp.Result)
	}
	return &resp, err
}

func call(client *http.Client, request *JSONRPC, response interface{}) (err error) {
	url := "https://kovan.infura.io/v3/fb298d4afd444cd5b7c5703b99d51f05"
	request.JSONRPC = "2.0"
	id++
	request.ID = id
	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(request)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("response: ", string(body))
	err = json.Unmarshal(body, &response)
	return
}
