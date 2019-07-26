// +build ignore

package main

import (
	"fmt"
	"math/big"

	"github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
)

func main() {
	h, err := sdk.EthBlockNumberInt64()
	if err != nil {
		fmt.Println("eth_blockNumber error: ", err)
		return
	}
	fmt.Println("height: ", h)

	// resp, err := sdk.EthGetTransactionReceiptResponse(txhash)
	// fmt.Println("Tx receipt: ", resp)
	resp, err := sdk.EthGetBalance("0xb0d2da0f43Cd2E44e4F3a38E24945F0ca0Ea95e2", "latest")
	if err != nil {
		fmt.Println("EthGetBalance error: ", err)
	}
	fmt.Println("response: ", resp)

	// txhash := "0x88efebd7469c323df3f7c1e57e51e586643d3df9c80ff18f70a48e93b13fea89"
	// txhash := "0xc07a33c29304e6eb92cfce2c9d95d916686648b973a86aa791efce6c35f98998"
	// txhash := "0xac9a4af995c2fb996379a8a95dd41d69a7a3d04335b971339aeb8419d0bb7db0"
	// txhash := "0xd1477743101cb2af7c712a3d1cf84790756b28ea74cf238c58f208926c08f77a"
	txhash := "0x9498cd360d1e295fe294a423265548af356bca1f49faf60560a27015702039b4"

	tx, err := sdk.EthGetTransactionByHash(txhash)
	fmt.Println("Tx: ", tx.GasPrice)

	receipt, err := sdk.EthGetTransactionReceipt(txhash)
	fmt.Println("Status: ", receipt.Status, "; ", receipt.Success(), "; ", receipt.GasUsed)

	v := mul(receipt.GasUsed, tx.GasPrice)
	fmt.Printf("mul: 0x%s ; %s\n", v.Text(16), v.Text(10))
}

func mul(gasUsed, gasPrice string) *big.Int {
	g := new(big.Int)
	g.SetString(gasUsed[2:], 16)
	gp := new(big.Int)
	gp.SetString(gasPrice[2:], 16)
	g = g.Mul(g, gp)
	// fee := fmt.Sprintf("0x%s\n", g.Text(16))
	// fmt.Println("fee: ", fee, "; ", g.Text(10))
	return g
}
