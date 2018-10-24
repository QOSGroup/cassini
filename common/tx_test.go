package common

import (
	"strings"
	"testing"

	motx "github.com/QOSGroup/cassini/mock/tx"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/stretchr/testify/assert"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestGetTxQcpHashCheck(t *testing.T) {
	tx := motx.NewTxQcpMock("abc", "xyz", 1, 99)
	event, err := Transform(tx)
	assert.NoError(t, err)

	hashStr := getHashStr(event, qcp.QcpHash)

	txo := motx.NewTxQcpMock("abc", "xyz", 1, 99)
	hashTxoStr := Bytes2HexStr(HashTx(txo))
	assert.Equal(t, hashStr, hashTxoStr)
}

func TestGetTxQcpHashCheckHeight(t *testing.T) {
	tx := motx.NewTxQcpMock("abc", "xyz", 1, 99)
	event, err := Transform(tx)
	assert.NoError(t, err)

	hashStr := getHashStr(event, qcp.QcpHash)

	txo := motx.NewTxQcpMock("abc", "xyz", 2, 99)
	hashTxoStr := Bytes2HexStr(HashTx(txo))
	assert.NotEqual(t, hashStr, hashTxoStr)
}

func TestGetTxQcpHashCheckFrom(t *testing.T) {
	tx := motx.NewTxQcpMock("abc", "xyz", 1, 99)
	event, err := Transform(tx)
	assert.NoError(t, err)

	hashStr := getHashStr(event, qcp.QcpHash)

	txo := motx.NewTxQcpMock("abcd", "xyz", 1, 99)
	hashTxoStr := Bytes2HexStr(HashTx(txo))
	assert.NotEqual(t, hashStr, hashTxoStr)
}

func TestGetTxQcpHashCheckTo(t *testing.T) {
	tx := motx.NewTxQcpMock("abc", "xyz", 1, 99)
	event, err := Transform(tx)
	assert.NoError(t, err)

	hashStr := getHashStr(event, qcp.QcpHash)

	txo := motx.NewTxQcpMock("abc", "axyz", 1, 99)
	hashTxoStr := Bytes2HexStr(HashTx(txo))
	assert.NotEqual(t, hashStr, hashTxoStr)
}

func TestGetTxQcpHashCheckSequence(t *testing.T) {
	tx := motx.NewTxQcpMock("abc", "xyz", 1, 99)
	event, err := Transform(tx)
	assert.NoError(t, err)

	hashStr := getHashStr(event, qcp.QcpHash)

	txo := motx.NewTxQcpMock("abc", "xyz", 1, 11)
	hashTxoStr := Bytes2HexStr(HashTx(txo))
	assert.NotEqual(t, hashStr, hashTxoStr)
}

func getHashStr(e *tmtypes.EventDataTx, key string) string {
	for _, kv := range e.Result.Tags {
		if strings.EqualFold(key, string(kv.Key)) {
			return Bytes2HexStr(kv.Value)
		}
	}
	return ""
}
