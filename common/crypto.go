package common

import (
	"encoding/hex"

	catypes "github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// Bytes2HexStr 字节数组转换为16 进制字符串
func Bytes2HexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HashTx 交易哈希算法，用以验证交易
func HashTx(tx *txs.TxQcp) []byte {
	return crypto.Sha256(tx.BuildSignatureBytes())
}

// UnmarshalKey returns the key ed25519.PrivKeyEd25519
func UnmarshalKey(base64key string) (*ed25519.PrivKeyEd25519, error) {
	caHex := "{\"type\": \"tendermint/PrivKeyEd25519\",\"value\": \"" + base64key + "\"}"
	var key ed25519.PrivKeyEd25519
	cdc := catypes.CreateCompleteCodec() //TODO
	err := cdc.UnmarshalJSON([]byte(caHex), &key)

	return &key, err
}
