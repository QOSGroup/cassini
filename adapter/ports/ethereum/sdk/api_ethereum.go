package sdk

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/QOSGroup/cassini/adapter/ports/ethereum/token"
)

// Token info on chain
type Token struct {
	Name     string `json:"name,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Decimals uint8  `json:"decimals,omitempty"`
}

// NewAccount create a new account on ethereum
func NewAccount(accountID, pass string) (*Account, error) {
	keydir := "/root/.ethereum/keystore"
	pass = "12345678"
	address, err := keystore.StoreKey(keydir, pass,
		keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	account := &Account{
		WalletAddress: address.Hex()}
	return account, nil
}

// QueryTokenInfo query token info on ethereum
func QueryTokenInfo(chain, tokenAddress string) (*Token, error) {
	client, err := ethclient.Dial(cli.url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer client.Close()
	address := common.HexToAddress(tokenAddress)
	instance, err := token.NewToken(address, client)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	t := &Token{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals}
	return t, nil
}
