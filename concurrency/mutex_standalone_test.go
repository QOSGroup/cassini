package concurrency

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	var m Mutex
	m = NewStandaloneMutex("abc")
	seq, err := m.Lock(0)
	assert.Error(t, err)
	assert.Equal(t, int64(1), seq)

	seq, err = m.Lock(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)

	var wg sync.WaitGroup

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			seq, err := m.Lock(1)
			assert.Equal(t, int64(2), seq)
			if err == nil {
				m.Unlock(true)
			}
			wg.Done()
		}()
	}
	m.Unlock(true)

	wg.Wait()

	wg.Add(1)
	go func() {
		seq, err := m.Lock(5)
		assert.Error(t, err)
		assert.Equal(t, int64(2), seq)
		wg.Done()
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		seq, err := m.Lock(2)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), seq)
		m.Unlock(true)
		wg.Done()
	}()
	wg.Wait()

	err = m.Update(15)
	assert.NoError(t, err)

	wg.Add(1)
	go func() {
		seq, err := m.Lock(15)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), seq)
		m.Unlock(true)
		wg.Done()
	}()
	wg.Wait()
}
