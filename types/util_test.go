package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//import (
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestBytesInt64(t *testing.T) {
//
//	bytes := Int64Bytes(256)
//	y, err := BytesInt64(bytes)
//
//	assert.NoError(t, err)
//	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0}, bytes)
//	assert.Equal(t, int64(256), y)
//}
//
//func TestGetMaxChainOutSequenceKey(t *testing.T) {
//
//	key := GetMaxChainOutSequenceKey("TEST")
//
//	assert.Equal(t, "sequence/out/TEST", key)
//}
//
//func TestGetChainOutTxsKey(t *testing.T) {
//
//	key := GetChainOutTxsKey("TEST", 111)
//
//	assert.Equal(t, "tx/out/TEST/111", key)
//}

func TestParseAddrs(t *testing.T) {
	p, as := ParseAddrs("etcd://127.0.0.1:8080,192.168.1.111:777")

	assert.Equal(t, p, "etcd")
	assert.Equal(t, len(as), 2)

	p, as = ParseAddrs("127.0.0.1:8080")

	assert.Equal(t, p, "")
	assert.Equal(t, len(as), 0)

	p, as = ParseAddrs("")

	assert.Equal(t, p, "")
	assert.Equal(t, len(as), 0)

	p, as = ParseAddrs("redis://192.168.1.111:8888,111.1.11.11:1111,222.2.222.222:2222")

	assert.Equal(t, p, "redis")
	assert.Equal(t, len(as), 3)
}
