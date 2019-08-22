package queue

import (
	"strings"
	"sync"

	"github.com/QOSGroup/cassini/commands"
	exporter "github.com/QOSGroup/cassini/prometheus"
	"github.com/spf13/viper"
)

var queues sync.Map

// NewProducer returns a new producer of message queue service
func NewProducer(subject string) (Producer, error) {
	queue := getQueue(subject)
	return queue.NewProducer()
}

// NewConsumer returns a new consumer of message queue service
func NewConsumer(subject string) (Consumer, error) {
	queue := getQueue(subject)
	return queue.NewConsumer()
}

func getQueue(subject string) Queue {
	var queue Queue
	q, loaded := queues.Load(subject)
	if !loaded {
		conf := viper.GetString(commands.FlagQueue)
		queue = newQueue(subject, conf)
		q, loaded = queues.LoadOrStore(subject, queue)
	}
	if loaded {
		queue = q.(Queue)
	}
	return queue
}

func newQueue(subject, conf string) Queue {
	if strings.HasPrefix(conf, "nats://") {
		return &NatsQueue{Subject: subject, Config: conf}
	}
	return &LocalQueue{Subject: subject, Config: conf}
}

// Queue define the queue of message queue service
type Queue interface {
	Init() error
	NewProducer() (Producer, error)
	NewConsumer() (Consumer, error)
}

// Producer define the producer of message queue service
type Producer interface {
	Subject() string
	Config() string
	Produce([]byte) error
}

// Listener for message listening
type Listener func([]byte, Consumer)

// Consumer define the consumer of message queue service
type Consumer interface {
	Subject() string
	Config() string
	Subscribe(Listener) error
}

// LocalQueue implements a inner version of message queue
type LocalQueue struct {
	Subject       string
	Config        string
	isInitialized bool
	ch            chan []byte
	sync.Mutex
}

// Init message queue
func (q *LocalQueue) Init() error {
	if q.isInitialized {
		return nil
	}
	q.Lock()
	defer q.Unlock()
	if q.isInitialized {
		return nil
	}
	queueSize := 100
	q.ch = make(chan []byte, queueSize)
	q.isInitialized = true
	exporter.Set(exporter.KeyQueueSize, float64(queueSize), "local")
	return nil
}

// NewProducer returns a new producer for the message queue
func (q *LocalQueue) NewProducer() (p Producer, err error) {
	q.Init()
	return &LocalProducer{queue: q, out: q.ch}, nil
}

// NewConsumer returns a new consumer for the message queue
func (q *LocalQueue) NewConsumer() (c Consumer, err error) {
	q.Init()
	return &LocalConsumer{queue: q, in: q.ch}, nil
}

// LocalProducer define the producer for local message queue based on channel
type LocalProducer struct {
	queue *LocalQueue
	out   chan<- []byte
}

// Produce sends data to local message queue
func (p *LocalProducer) Produce(data []byte) error {
	p.out <- data
	return nil
}

// Subject returns subject of local message queue
func (p *LocalProducer) Subject() string {
	return p.queue.Subject
}

// Config returns config of local message queue
func (p *LocalProducer) Config() string {
	return p.queue.Config
}

// LocalConsumer define the comsumer for local message queue based on channel
type LocalConsumer struct {
	queue *LocalQueue
	in    <-chan []byte
}

// Subscribe sets the listener for local message queue based on channel
func (c *LocalConsumer) Subscribe(listener Listener) error {
	go func() {
		for {
			data := <-c.in
			listener(data, c)
		}
	}()
	return nil
}

// Subject returns subject of local message queue
func (c *LocalConsumer) Subject() string {
	return c.queue.Subject
}

// Config returns config of local message queue
func (c *LocalConsumer) Config() string {
	return c.queue.Config
}
