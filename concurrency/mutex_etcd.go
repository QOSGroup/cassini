package concurrency

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	v3 "github.com/coreos/etcd/clientv3"
	v3c "github.com/coreos/etcd/clientv3/concurrency"
)

// NewEtcdMutex new a mutex for a etcd implementation.
func NewEtcdMutex(chainID string, addrs []string) (*EtcdMutex, error) {
	cli, err := v3.New(v3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Errorf("New client error: %s", err)
		return nil, err
	}

	var sess *v3c.Session
	sess, err = v3c.NewSession(cli, v3c.WithTTL(5))
	if err != nil {
		cli.Close()
		log.Errorf("New session error: %s", err)
		return nil, err
	}

	m := &EtcdMutex{chainID: chainID,
		client:  cli,
		session: sess}

	return m, nil
}

// EtcdMutex implements a distributed lock based on etcd.
type EtcdMutex struct {
	chainID  string
	client   *v3.Client
	session  *v3c.Session
	mutex    *v3c.Mutex
	sequence int64
	locked   bool
}

// Lock get lock
func (e *EtcdMutex) Lock(sequence int64) (int64, error) {
	e.mutex = v3c.NewMutex(e.session, e.chainID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()

	var err error
	err = e.mutex.Lock(ctx)
	seq := e.get()
	if err != nil {
		log.Errorf("Lock error: %s", err)
		return seq, err
	}
	e.locked = true

	if sequence != seq {
		defer func() {
			err := e.Unlock(false)
			if err != nil {
				log.Error("Unlock error: ", err)
			}
		}()
		err = fmt.Errorf("Wrong sequence(%d), current sequence(%d) in lock",
			sequence, seq)
		log.Error(err)
		return seq, err
	}
	e.sequence = sequence

	log.Debugf("Get lock success, %s: %d", e.chainID, e.sequence)
	return e.sequence, nil
}

// Update update the lock
func (e *EtcdMutex) Update(sequence int64) error {
	mux := v3c.NewMutex(e.session, e.chainID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()

	var err error
	err = mux.Lock(ctx)
	if err != nil {
		log.Errorf("Update lock sequence(%d) error: %s", sequence, err)
		return err
	}
	defer func() {
		mux.Unlock(ctx)
	}()
	seq := e.get()
	if sequence > seq {
		err = e.put(sequence)
		if err != nil {
			log.Errorf("Update sequence(%d), current sequence(%d) in lock, error: %s",
				sequence, seq, err)
			return err
		}
	}
	e.sequence = sequence
	log.Debugf("Upadte lock success, %s: %d", e.chainID, e.sequence)
	return nil
}

// Unlock unlock the lock
func (e *EtcdMutex) Unlock(success bool) (err error) {
	if !e.locked {
		return nil
	}
	if success {
		e.sequence++
		err = e.put(e.sequence)
		if err != nil {
			log.Errorf("Put key value error when unlock: ", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	err = e.mutex.Unlock(ctx)
	if err != nil {
		log.Errorf("Unlock error: ", err)
		return
	}
	e.locked = false
	log.Debugf("Unlock success, %s: %d", e.chainID, e.sequence)
	return
}

// Close close the lock
func (e *EtcdMutex) Close() error {
	e.session.Close()
	return e.client.Close()
}

func (e *EtcdMutex) get() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	// var resp *v3.GetResponse
	resp, err := e.client.Get(ctx, e.chainID)
	if err != nil {
		log.Error("Get key value error: ", err)
		return -1
	}
	for _, kv := range resp.Kvs {
		if strings.EqualFold(string(kv.Key), e.chainID) {
			var seq int64
			seq, err = types.ParseSequence(kv.Value)
			if err != nil {
				log.Error("Parse key value error: ", err)
				return -1
			}
			return seq
		}
	}
	return -1
}

func (e *EtcdMutex) put(sequence int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	// var resp *v3.PutResponse
	_, err := e.client.Put(ctx, e.chainID, fmt.Sprintf("%d", sequence))
	if err != nil {
		log.Error("Put key value error: ", err)
	}
	return err
}
