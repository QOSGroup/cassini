package sdk

import (
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
)

func newKeyFromECDSA(priECDSA *ecdsa.PrivateKey) *keystore.Key {
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(priECDSA.PublicKey),
		PrivateKey: priECDSA,
	}
	return key
}

// newKey create a new key for ethereum account
func newKey(rand io.Reader) (*keystore.Key, error) {
	priECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(priECDSA), nil
}
