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
	err := RegisterAdapter(ip, port, chain)
	err = RegisterAdapter(ip, port+1, chain)
	err = RegisterAdapter(ip, port+2, chain)

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
	a := &qosAdapter{
		chain: "test",
		ip:    "192.168.1.111",
		port:  26657}
	key := GetAdapterKey(a)
	assert.Equal(t, "test://192.168.1.111:26657", key)

	a = &qosAdapter{
		chain: "target-chain",
		ip:    "127.0.0.1",
		port:  8080}
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
