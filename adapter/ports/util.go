package ports

import "fmt"

// GetAdapterKey Gen an adapter key
func GetAdapterKey(a Adapter) string {
	return fmt.Sprintf("%s://%s:%d", a.GetChain(), a.GetIP(), a.GetPort())
}

// GetNodeAddress Gen a node address
func GetNodeAddress(a Adapter) string {
	return fmt.Sprintf("%s:%d", a.GetIP(), a.GetPort())
}

// Consensus2of3 Calculate number of consensus
func Consensus2of3(value int) int {
	return (value*2 + 2) / 3
}
