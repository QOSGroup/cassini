package ports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFerry(t *testing.T) {
	// mc := config.TestQscMockConfig()
	// cancelMock, err := mock.StartMock(*mc)
	// defer cancelMock()
	// addr := mc.RPC.NodeAddress
	// ipPort := strings.SplitN(addr, ":", 2)
	// if len(ipPort) != 2 {
	// 	err = fmt.Errorf("Ip and port parse error: %v", addr)
	// 	assert.NoError(t, err)
	// }
	ip := "127.0.0.1"
	port := 27657
	chain := "qos"
	conf := &AdapterConfig{
		ChainName: chain,
		IP:        ip,
		Port:      port}
	err := RegisterAdapter(conf)
	conf.Port++
	err = RegisterAdapter(conf)
	conf.Port++
	err = RegisterAdapter(conf)

	assert.NoError(t, err)

	var ads map[string]Adapter
	ads, err = GetAdapters(chain)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(ads))

	for _, a := range ads {
		total, consensusNumber := a.Count()
		assert.Equal(t, 3, total)
		assert.Equal(t, 2, consensusNumber)
	}
}

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
}
