// +build ignore

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	// store "./contracts" // for demo
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("wss://kovan.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xEAd5C972Fe8Bbf6f725Ab8A4C7E9d40E15f35241")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(10705633),
		ToBlock:   big.NewInt(10710975),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, vLog := range logs {
		fmt.Println("block hash: ", vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
		fmt.Println("block number: ", vLog.BlockNumber)   // 2394201
		fmt.Println("tx hash: ", vLog.TxHash.Hex())       // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6
		fmt.Println("address: ", vLog.Address.Hex())
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		// err := contractAbi.Unpack(&event, "ItemSet", vLog.Data)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		fmt.Println("data: ", hex.EncodeToString(vLog.Data))

		fmt.Println("key: ", string(event.Key[:]))     // foo
		fmt.Println("value: ", string(event.Value[:])) // bar

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
			fmt.Printf("topic[%d]: %s\n", i, topics[i])
		}

		// 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
	}

	// event Transfer erc20
	eventSignature := []byte("Transfer(address,address,uint256)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println("event hash: ", hash.Hex()) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4

	eventSignature = []byte("transfer(address,uint256)")
	hash = crypto.Keccak256Hash(eventSignature)
	fmt.Println("method id: ", hash.Hex()[0:10])

}
