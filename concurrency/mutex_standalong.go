package concurrency

import (
	"fmt"
	"sync"

	"github.com/QOSGroup/cassini/config"
)

// NewStandaloneMutex new a mutex for a standalone implementation.
func NewStandaloneMutex(conf *config.QscConfig) *StandaloneMutex {
	return &StandaloneMutex{
		chainID:  conf.Name,
		sequence: 1}
}

// StandaloneMutex implements a standalone version for single process.
type StandaloneMutex struct {
	chainID  string
	sequence int64
	mux      sync.Mutex
}

// Lock get lock
func (s *StandaloneMutex) Lock(sequence int64) (int64, error) {
	if sequence < s.sequence {
		return s.sequence, fmt.Errorf("Wrong sequence(%d): lock sequence(%d)",
			sequence, s.sequence)
	}
	s.mux.Lock()
	if sequence < s.sequence {
		s.mux.Unlock()
		return s.sequence, fmt.Errorf("Wrong sequence(%d): lock sequence(%d)",
			sequence, s.sequence)
	}
	s.sequence = sequence
	return s.sequence, nil
}

// Unlock unlock the lock
func (s *StandaloneMutex) Unlock(success bool) error {
	if success {
		s.sequence++
	}
	s.mux.Unlock()
	return nil
}
