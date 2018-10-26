package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesInt64(t *testing.T) {

	bytes := Int64Bytes(256)
	y, err := BytesInt64(bytes)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0}, bytes)
	assert.Equal(t, int64(256), y)
}

func TestGetMaxChainOutSequenceKey(t *testing.T) {

	key := GetMaxChainOutSequenceKey("TEST")

	assert.Equal(t, "sequence/out/TEST", key)
}

func TestGetChainOutTxsKey(t *testing.T) {

	key := GetChainOutTxsKey("TEST", 111)

	assert.Equal(t, "tx/out/TEST/111", key)
}
