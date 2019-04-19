// +build ignore

package main

import (
	"fmt"

	"github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
)

func main() {
	h, err := sdk.EthBlockNumberInt64()
	if err != nil {
		fmt.Println("eth_blockNumber error: ", err)
		return
	}
	fmt.Println("height: ", h)

	txhash := "0x88efebd7469c323df3f7c1e57e51e586643d3df9c80ff18f70a48e93b13fea89"
	// txhash := "0x4e9255e66cd7a948d600d87ca1fb3e8dc2ec3edfb72659e12d0918177812dae4"
	resp, err := sdk.EthGetTransactionReceiptResponse(txhash)
	fmt.Println("Tx receipt: ", resp)

	receipt, err := sdk.EthGetTransactionReceipt(txhash)
	fmt.Println("Status: ", receipt.Status, "; ", receipt.Success())

	resp, err = sdk.EthGetBalance("0xb0d2da0f43Cd2E44e4F3a38E24945F0ca0Ea95e2", "latest")
	if err != nil {
		fmt.Println("EthGetBalance error: ", err)
	}
	fmt.Println("response: ", resp)
}
