// Package concurrency provides support for distributed concurrency capabilities
//
// The goal is to simplify the development of concurrency capabilities.
package concurrency

import (
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/types"
)

// Mutex distributed lock interface
type Mutex interface {
	// Lock get lock through chain id and sequence
	//
	// If it returned an error, indicates that the call failed.
	// Whether successful or not,
	// the current sequence saved in the distributed lock is returned.
	// Negative sequence(<0) are returned unless there are some unknown exceptions.
	Lock(sequence int64) (int64, error)

	// Unlock after successfully acquiring the lock, the lock needs to be unlocked.
	//
	// If it returned an error, indicates that the call failed.
	Unlock(success bool) error
}

// NewMutex new mutex based on configuration.
func NewMutex(conf *config.QscConfig) Mutex {

	protocol, addrs := types.ParseAddrs(conf.Lock)

	var m Mutex
	switch protocol {
	case "etcd":
		m = NewEtcdMutex(protocol, addrs, conf)
	default:
		m = NewStandaloneMutex(conf)
	}
	return m
}
