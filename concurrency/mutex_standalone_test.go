package concurrency

import (
	"sync"
	"testing"

	"github.com/QOSGroup/cassini/config"
	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	c := &config.QscConfig{Name: "abc"}
	var m Mutex
	m = NewStandaloneMutex(c)
	seq, err := m.Lock(0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), seq)

	var wg sync.WaitGroup

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			seq, err = m.Lock(0)
			assert.Equal(t, int64(1), seq)
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
		seq, err = m.Lock(5)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), seq)
		m.Unlock(true)
		wg.Done()
	}()
	wg.Wait()

	seq, err = m.Lock(3)
	assert.Error(t, err)
	assert.Equal(t, int64(6), seq)
}
