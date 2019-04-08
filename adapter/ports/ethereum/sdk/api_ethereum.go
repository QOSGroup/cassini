package sdk

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

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
