package queue

import (
	"fmt"
	"sync"
	"testing"

	"github.com/QOSGroup/cassini/commands"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_getLocalQueue(t *testing.T) {
	viper.Set(commands.FlagQueue, "local")

	q := getQueue("test")

	q2 := getQueue("test")

	assert.Equal(t, q, q2)
}

var wg sync.WaitGroup

func Test_NewProducer(t *testing.T) {
	viper.Set(commands.FlagQueue, "local")

	p, err := NewProducer("test")
	assert.NoError(t, err)

	assert.Equal(t, "local", p.Config(), "get wrong producer")

	if p != nil {
		p.Produce([]byte("test"))
		wg.Add(1)
	}
}

func Test_NewConsumer(t *testing.T) {
	viper.Set(commands.FlagQueue, "local")

	c, err := NewConsumer("test")
	assert.NoError(t, err)

	if c != nil {
		c.Subscribe(func(data []byte, consumer Consumer) {
			t.Logf("queue %s get: %s", consumer.Subject(), string(data))
			wg.Done()
		})
	}
	wg.Wait()
}

func Test_Subscribe(t *testing.T) {
	viper.Set(commands.FlagQueue, "local")

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

	c, err := NewConsumer("test2")
	assert.NoError(t, err)
	t.Log("waiting")
	if c != nil {
		c.Subscribe(func(data []byte, consumer Consumer) {
			t.Logf("queue %s get: %s", consumer.Subject(), string(data))
			wg.Done()
			// t.Log("done")
		})
	}
	wg.Wait()
}

// func Benchmark_LocalQueue(b *testing.B) {
// }

func Benchmark_Parallel_LocalQueue(b *testing.B) {
	viper.Set(commands.FlagQueue, "local")

	b.ReportAllocs()
	var counter int
	c, err2 := NewConsumer("test_parallel")
	if err2 == nil && c != nil {
		c.Subscribe(func(data []byte, consumer Consumer) {
			counter++
		})
	}

	for i := 0; i < 5; i++ {
		b.RunParallel(func(pb *testing.PB) {
			p, err := NewProducer("test_parallel")
			if err == nil {
				if p != nil {
					i := 0
					for pb.Next() {
						i++
						p.Produce([]byte(fmt.Sprintf("testing-%d", i)))
					}
				}
			}
		})
	}
	b.Log("counter: ", counter)
}
