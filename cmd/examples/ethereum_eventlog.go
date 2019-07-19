// +build ignore

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/QOSGroup/cassini/adapter/ports/ethereum/token"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

func main() {
	client, err := ethclient.Dial("wss://kovan.infura.io/ws")
	// client, err := ethclient.Dial("wss://mainnet.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xEAd5C972Fe8Bbf6f725Ab8A4C7E9d40E15f35241")
	// contractAddress := common.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(10705633),
		ToBlock:   big.NewInt(10710975),
		// FromBlock: big.NewInt(7534084),
		// ToBlock:   big.NewInt(7535084),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	// event Transfer erc20
	eventSignature := []byte("Transfer(address,address,uint256)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println("event hash: ", hash.Hex())

	eventSignature = []byte("transfer(address,uint256)")
	hash = crypto.Keccak256Hash(eventSignature)
	fmt.Println("method id: ", hash.Hex()[0:10], hash.Hex())

	for _, vLog := range logs {
		fmt.Println("block hash: ", vLog.BlockHash.Hex())
		fmt.Println("block number: ", vLog.BlockNumber)
		fmt.Println("tx hash: ", vLog.TxHash.Hex())
		fmt.Println("address: ", vLog.Address.Hex())
		fmt.Println("data: ", hex.EncodeToString(vLog.Data))

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
			fmt.Printf("topic[%d]: %s\n", i, topics[i])
		}

		var event EventTransfer

		event.From = common.HexToAddress(topics[1])
		event.To = common.HexToAddress(topics[2])

		err := contractAbi.Unpack(&event, "Transfer", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("From: %s\n", event.From.Hex())
		fmt.Printf("To: %s\n", event.To.Hex())
		fmt.Printf("Tokens: %s\n", event.Tokens.String())

	}

}
