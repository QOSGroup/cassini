package sdk

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// NewAccount create a new account on ethereum
func NewAccount(accountID, pass string) (*Account, error) {
	key, err := newKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	keyjson, err := keystore.EncryptKey(key, pass,
		keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	account := &Account{
		WalletAddress: key.Address.Hex(),
		EncryptedKey:  string(keyjson)}
	return account, nil
}
