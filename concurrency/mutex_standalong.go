package concurrency

import (
	"fmt"
	"sync"
)

// NewStandaloneMutex new a mutex for a standalone implementation.
func NewStandaloneMutex(name string) *StandaloneMutex {
	return &StandaloneMutex{
		chainID:  name,
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

// Close close the lock
func (s *StandaloneMutex) Close() error {
	return nil
}
