package concurrency

import "github.com/QOSGroup/cassini/config"

// NewEtcdMutex new a mutex for a etcd implementation.
func NewEtcdMutex(protocol string, addrs []string, conf *config.QscConfig) *EtcdMutex {

	return nil
}

// EtcdMutex implements a distributed lock based on etcd.
type EtcdMutex struct {
}
