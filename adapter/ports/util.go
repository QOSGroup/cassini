package ports

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/QOSGroup/cassini/log"
)

// GetAdapterKey Gen an adapter key
func GetAdapterKey(a Adapter) string {
	return genAdapterKey(a.GetChainName(), a.GetIP(), a.GetPort())
}

// GetAdapterKeyByConfig Gen an adapter key from adapter config
func GetAdapterKeyByConfig(c *AdapterConfig) string {
	return genAdapterKey(c.ChainName, c.IP, c.Port)
}

func genAdapterKey(chainName, ip string, port int) string {
	return fmt.Sprintf("%s://%s:%d", chainName, ip, port)
}

// GetNodeAddress Gen a node address
func GetNodeAddress(a Adapter) string {
	return fmt.Sprintf("%s:%d", a.GetIP(), a.GetPort())
}

// Consensus2of3 Calculate number of consensus
func Consensus2of3(value int) int {
	return (value*2 + 2) / 3
}

// ParseNodeAddress parse node address to get ip and port
func ParseNodeAddress(nodeAddr string) (string, int, error) {
	addrs := strings.Split(nodeAddr, ":")
	var msg string
	if len(addrs) == 2 {
		port, err := strconv.Atoi(addrs[1])
		if err == nil {
			return addrs[0], port, nil
		}
		msg = fmt.Sprintf("Node address parse error: %s, %v", nodeAddr, err)
	} else {
		msg = fmt.Sprintf("Can not parse node address %s", nodeAddr)
	}
	log.Errorf(msg)
	return "", -1, errors.New(msg)
}
