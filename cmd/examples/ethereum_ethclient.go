package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/log"
)

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
	height := int64(10705633)
	h := int64(0)
	var err error
	for true {
		h, err = sdk.EthBlockNumberInt64()
		if err != nil {
			fmt.Println("eth_blockNumber error: ", err)
			continue
		} else if height == h {
			continue
		}
		if height < h {
			height++
		} else {
			height = h
		}
		fmt.Printf("eth_blockNumber %d\n", height)

		protocol, err := sdk.EthProtocolVersionResponse()
		if err != nil {
			fmt.Println("eth_protocolVersion error: ", err)
		} else {
			fmt.Println("eth_protocolVersion resp: ", protocol)
		}
		var result *sdk.ResultBlock
		for result == nil {
			result, err = sdk.EthGetBlockByNumber(height)
			if err != nil {
				fmt.Println("eth_getBlockByNumber error: ", err)
			}
			if result == nil {
				fmt.Printf("block %d is nil\n", height)
			}
		}
		checkBlock(result)

		// result, err = sdk.EthGetBlockByNumber(10686098)
		// if err != nil {
		// 	fmt.Println("eth_getBlockByNumber error: ", err)
		// } else {
		// 	checkBlock(result)
		// }
		time.Sleep(1 * time.Second)
	}
}

func checkBlock(result *sdk.ResultBlock) {
	fmt.Println("eth_getBlockByNumber resp: ",
		result.Difficulty, len(result.Transactions))
	for i, tx := range result.Transactions {
		value, err := strconv.ParseInt(tx.Value[2:], 16, 64)
		if err != nil {
			log.Errorf("value: %s parse error: %v", tx.Value, err)
			continue
		}
		fmt.Println("tx in block: ", i, "; ",
			tx.From, " -> ", tx.To, " : ", value)
		if strings.EqualFold(tx.To,
			// "0x3d947eB8c366D2416468675cEDd00fd311D70dFB") {
			"0xb0d2da0f43Cd2E44e4F3a38E24945F0ca0Ea95e2") {
			fmt.Println("check address: ",
				tx.To, "; ", tx.TransactionIndex,
				" value: ", value)
		}
	}
}
