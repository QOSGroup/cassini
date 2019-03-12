package sdk

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// NewAccount create a new account on ethereum
func NewAccount(name, password string) (*keystore.Key, error) {
	// conf = Config()
	// scryptN, scryptP, _, err := conf.Node.AccountConfig()
	// if err != nil {
	// 	utils.Fatalf("Configuration error: %v", err)
	// }
	return NewKey(rand.Reader)
}
