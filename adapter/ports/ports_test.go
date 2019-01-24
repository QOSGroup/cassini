package ports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdapter(t *testing.T) {
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

func TestRegisterAdapter(t *testing.T) {
	chainName := "cassini-test"
	var testBuilder Builder = func(config AdapterConfig) (AdapterService, error) {
		a := &QosAdapter{config: &config}
		return a, nil
	}
	GetPortsIncetance().RegisterBuilder(chainName, testBuilder)
	ip := "192.168.1.100"
	port := 9999
	conf := &AdapterConfig{
		ChainName: chainName,
		IP:        ip,
		Port:      port}
	err := RegisterAdapter(conf)
	assert.NoError(t, err)

	c := GetPortsIncetance().Count(chainName)
	assert.Equal(t, 1, c)

	err = RegisterAdapter(conf)
	assert.Error(t, err)
	err = RegisterAdapter(conf)
	assert.Error(t, err)

	conf.Port++
	err = RegisterAdapter(conf)
	assert.NoError(t, err)
	conf.Port++
	err = RegisterAdapter(conf)
	assert.NoError(t, err)

	var ads map[string]Adapter
	ads, err = GetPortsIncetance().Get(chainName)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ads))
}
