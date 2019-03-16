package sdk

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common/math"
)

// NewAccount create a new account on ethereum
func NewAccount(name, password string) (*Account, error) {
	key, err := newKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	bytes := math.PaddedBigBytes(key.PrivateKey.D, 32)
	account := &Account{
		WalletAddress: key.Address.Hex(),
		PrivateKey:    hex.EncodeToString(bytes)}
	return account, nil
}
