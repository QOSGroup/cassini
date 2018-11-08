package concurrency

import "github.com/QOSGroup/cassini/config"

// NewEtcdMutex new a mutex for a etcd implementation.
func NewEtcdMutex(protocol string, addrs []string, conf *config.QscConfig) *EtcdMutex {

	return nil
}

// EtcdMutex implements a distributed lock based on etcd.
type EtcdMutex struct {
}

// Lock get lock
func (e *EtcdMutex) Lock(sequence int64) (int64, error) {

	return 1, nil
}

// Unlock unlock the lock
func (e *EtcdMutex) Unlock(success bool) error {

	return nil
}
