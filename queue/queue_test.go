package queue

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getQueue(t *testing.T) {
	q := getQueue("test")

	q2 := getQueue("test")

	assert.Equal(t, q, q2)
}

func Test_NewProducer(t *testing.T) {
	p, err := NewProducer("test")
	assert.NoError(t, err)

	if p != nil {
		p.Produce([]byte("test"))
	}
}

func Test_NewComsumer(t *testing.T) {
	c, err := NewComsumer("test")
	assert.NoError(t, err)

	if c != nil {
		c.Subscribe(func(data []byte) {
			t.Logf("get: %s", string(data))
		})
	}
}

func Test_Subscribe(t *testing.T) {
	var wg sync.WaitGroup

	p, err := NewProducer("test2")
	assert.NoError(t, err)
	if p != nil {
		for i := 0; i < 3; i++ {
			p.Produce([]byte(fmt.Sprintf("testing_%d", i)))
			wg.Add(1)
			// t.Log("add")
		}
	}

	c, err := NewComsumer("test2")
	assert.NoError(t, err)
	t.Log("waiting")
	if c != nil {
		c.Subscribe(func(data []byte) {
			t.Logf("get: %s", string(data))
			wg.Done()
			// t.Log("done")
		})
	}
	wg.Wait()
}
