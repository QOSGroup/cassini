package queue

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/viper"

	"github.com/QOSGroup/cassini/commands"
	"github.com/stretchr/testify/assert"
)

func Test_getNatsQueue(t *testing.T) {
	viper.Set(commands.FlagQueue, "nats://127.0.0.1:4222")

	q := getQueue("nats-test")

	q2 := getQueue("nats-test")

	assert.Equal(t, q, q2)
}

var wgNats sync.WaitGroup

func Test_NatsNewConsumer(t *testing.T) {
	viper.Set(commands.FlagQueue, "nats://127.0.0.1:4222")

	c, err := NewConsumer("nats-test1")
	assert.NoError(t, err)

	if c != nil {
		c.Subscribe(func(data []byte, consumer Consumer) {
			t.Logf("queue %s get: %s", c.Subject(), string(data))
			wgNats.Done()
		})
	}
	wgNats.Wait()
}

func Test_NatsNewProducer(t *testing.T) {
	viper.Set(commands.FlagQueue, "nats://127.0.0.1:4222")

	p, err := NewProducer("nats-test1")
	assert.NoError(t, err)

	assert.Equal(t, "nats://127.0.0.1:4222", p.Config(), "get wrong producer")

	if p != nil {
		p.Produce([]byte("test msg"))
		wgNats.Add(1)
	}
}

func Benchmark_NatsProducer(b *testing.B) {
	viper.Set(commands.FlagQueue, "nats://127.0.0.1:4222")

	p, err := NewProducer("Benchmark_NatsProducer")
	if err != nil {
		b.Errorf("connect to nats error: %v", err)
	}

	if !strings.EqualFold("nats://127.0.0.1:4222", p.Config()) {
		b.Errorf("get wrong producer: %s", p.Config())
	}

	for i := 0; i < b.N; i++ { //1000	   1.841105 ms/op
		p.Produce([]byte(fmt.Sprintf("testing-%d", i)))
	}
}

func Benchmark_NatsConsumer(b *testing.B) {
	viper.Set(commands.FlagQueue, "nats://127.0.0.1:4222")

	c, err := NewConsumer("Benchmark_NatsConsumer")
	if err != nil {
		b.Errorf("connect to nats error: %v", err)
	}

	DEFAULTMSG := "message for Benchmark_NatsConsumer"
	i := 0
	listener := func(data []byte, _ Consumer) {
		i++
		// log.Infof("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
		if !strings.EqualFold(string(data), DEFAULTMSG) {
			b.Error("get wrong message")
		}
	}
	c.Subscribe(listener)

	p, err := NewProducer("Benchmark_NatsConsumer")
	if err != nil {
		b.Errorf("connect to nats error: %v", err)
	}
	for i := 0; i < b.N; i++ { //30000	     51369 ns/op
		p.Produce([]byte(DEFAULTMSG))
	}
	b.Logf("consumer message: %d", i)
}

func Benchmark_Parallel_NatsQueue(b *testing.B) {
	viper.Set(commands.FlagQueue, "nats://127.0.0.1:4222")

	b.ReportAllocs()
	var counter int
	c, err2 := NewConsumer("nats-test_parallel")
	if err2 == nil && c != nil {
		c.Subscribe(func(data []byte, consumer Consumer) {
			counter++
		})
	}

	for i := 0; i < 5; i++ {
		b.RunParallel(func(pb *testing.PB) {
			p, err := NewProducer("nats-test_parallel")
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
