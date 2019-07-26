package queue

import (
	"strings"
	"sync"

	"github.com/QOSGroup/cassini/commands"
	"github.com/spf13/viper"
)

var queues sync.Map

// NewProducer returns a new producer of message queue service
func NewProducer(subject string) (Producer, error) {
	queue := getQueue(subject)
	return queue.NewProducer()
}

// NewComsumer returns a new comsumer of message queue service
func NewComsumer(subject string) (Comsumer, error) {
	queue := getQueue(subject)
	return queue.NewComsumer()
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
	NewComsumer() (Comsumer, error)
}

// Producer define the producer of message queue service
type Producer interface {
	Subject() string
	Config() string
	Produce([]byte) error
}

// Listener for message listening
type Listener func([]byte, Comsumer)

// Comsumer define the comsumer of message queue service
type Comsumer interface {
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
	q.ch = make(chan []byte, 100)
	q.isInitialized = true
	return nil
}

// NewProducer returns a new producer for the message queue
func (q *LocalQueue) NewProducer() (p Producer, err error) {
	q.Init()
	return &LocalProducer{queue: q, out: q.ch}, nil
}

// NewComsumer returns a new comsumer for the message queue
func (q *LocalQueue) NewComsumer() (c Comsumer, err error) {
	q.Init()
	return &LocalComsumer{queue: q, in: q.ch}, nil
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

// LocalComsumer define the comsumer for local message queue based on channel
type LocalComsumer struct {
	queue *LocalQueue
	in    <-chan []byte
}

// Subscribe sets the listener for local message queue based on channel
func (c *LocalComsumer) Subscribe(listener Listener) error {
	go func() {
		for {
			data := <-c.in
			listener(data, c)
		}
	}()
	return nil
}

// Subject returns subject of local message queue
func (c *LocalComsumer) Subject() string {
	return c.queue.Subject
}

// Config returns config of local message queue
func (c *LocalComsumer) Config() string {
	return c.queue.Config
}
