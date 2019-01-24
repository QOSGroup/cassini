package ports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAdapterKey(t *testing.T) {
	conf := &AdapterConfig{
		ChainName: "test",
		IP:        "192.168.1.111",
		Port:      26657}
	a := &QosAdapter{config: conf}
	key := GetAdapterKey(a)
	assert.Equal(t, "test://192.168.1.111:26657", key)

	conf = &AdapterConfig{
		ChainName: "target-chain",
		IP:        "127.0.0.1",
		Port:      8080}
	a = &QosAdapter{config: conf}
	key = GetAdapterKey(a)
	assert.Equal(t, "target-chain://127.0.0.1:8080", key)
}

func TestGetAdapterKeyByConfig(t *testing.T) {
	conf := &AdapterConfig{
		ChainName: "test",
		IP:        "192.168.1.111",
		Port:      26657}
	key := GetAdapterKeyByConfig(conf)
	assert.Equal(t, "test://192.168.1.111:26657", key)

	conf = &AdapterConfig{
		ChainName: "target-chain",
		IP:        "127.0.0.1",
		Port:      8080}
	key = GetAdapterKeyByConfig(conf)
	assert.Equal(t, "target-chain://127.0.0.1:8080", key)
}

func TestConsensus2of3(t *testing.T) {
	c := Consensus2of3(3)
	assert.Equal(t, 2, c)

	c = Consensus2of3(4)
	assert.Equal(t, 3, c)

	c = Consensus2of3(5)
	assert.Equal(t, 4, c)
}

func TestParseNodeAddress(t *testing.T) {
	ip, port, err := ParseNodeAddress("192.168.1.111:26657")
	assert.NoError(t, err)
	assert.Equal(t, 26657, port)
	assert.Equal(t, "192.168.1.111", ip)

	ip, port, err = ParseNodeAddress("192.168.1.111:a")
	assert.Error(t, err)

	ip, port, err = ParseNodeAddress("123,123,123")
	assert.Error(t, err)
}
