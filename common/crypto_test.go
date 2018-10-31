package common

import (
	"encoding/hex"
	"testing"

	mtx "github.com/QOSGroup/cassini/mock/tx"
	"github.com/stretchr/testify/assert"
)

func TestBytes2HexStr(t *testing.T) {
	bytes, err := hex.DecodeString("ab")

	assert.NoError(t, err)

	assert.Equal(t, "ab", Bytes2HexStr(bytes))
}

func TestHashTx(t *testing.T) {
	hash := "08291d0b3d1e84f91f656574f0366d0c2de962288f33c417e81bb0ebdd588d21"
	tx := mtx.NewTxQcpMock("qsc", "qos", 1, 99)
	assert.Equal(t, hash, Bytes2HexStr(HashTx(tx)))

	hash = "8bdc305f444399f0a2a3c4a8858ab494b8192dc11a7768033b18eb9530dc4a76"
	tx = mtx.NewTxQcpMock("qos", "qsc", 19, 111)
	assert.Equal(t, hash, Bytes2HexStr(HashTx(tx)))
}
