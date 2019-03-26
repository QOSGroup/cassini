package main

import (
	"fmt"
	"os"
	"time"

	"github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
)

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

func main() {
	go func() {
		callEthereum()
	}()

	common.KeepRunning(func(sig os.Signal) {
		sdk.Close()
		log.Info("cancel done.")
	})

}

func callEthereum() {
	for true {
		height, err := sdk.EthBlockNumber()
		if err != nil {
			fmt.Println("eth_blockNumber error: ", err)
		} else {
			fmt.Printf("eth_blockNumber %d\n", height)
		}

		protocol, err := sdk.EthProtocolVersion()
		if err != nil {
			fmt.Println("eth_protocolVersion error: ", err)
		} else {
			fmt.Println("eth_protocolVersion resp: ", protocol)
		}

		result, err := sdk.EthGetBlockByNumber(height)
		if err != nil {
			fmt.Println("eth_getBlockByNumber error: ", err)
		} else {
			fmt.Println("eth_getBlockByNumber resp: ",
				result.Difficulty, len(result.Transactions))
			for i, tx := range result.Transactions {
				fmt.Println("tx in block: ", i, "; ",
					tx.From, " -> ", tx.To)
			}
		}
		time.Sleep(2 * time.Second)
	}
}
