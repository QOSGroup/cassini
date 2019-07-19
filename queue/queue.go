package queue

import (
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
	return &LocalQueue{Subject: subject}
}

// Queue define the queue of message queue service
type Queue interface {
	Init()
	NewProducer() (Producer, error)
	NewComsumer() (Comsumer, error)
}

// Producer define the producer of message queue service
type Producer interface {
	Produce([]byte) error
}

// Listener for message listening
type Listener func(string, []byte)

// Comsumer define the comsumer of message queue service
type Comsumer interface {
	Subscribe(Listener) error
}

// LocalQueue implements a inner version of message queue
type LocalQueue struct {
	Subject       string
	isInitialized bool
	sync.Mutex
	ch chan []byte
}

// Init message queue
func (q *LocalQueue) Init() {
	if q.isInitialized {
		return
	}
	q.Lock()
	defer q.Unlock()
	if q.isInitialized {
		return
	}
	q.ch = make(chan []byte, 100)
	q.isInitialized = true
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
			listener(c.queue.Subject, data)
		}
	}()
	return nil
}
