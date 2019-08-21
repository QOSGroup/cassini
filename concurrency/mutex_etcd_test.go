package concurrency

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/QOSGroup/cassini/config"
	"github.com/etcd-io/etcd/embed"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestEtcdMutex(t *testing.T) {
	confStr := `
embedEtcd:	true
useEtcd:	true
lock:		etcd://127.0.0.1:2379
etcd:
  name:				test-cassini
  advertise:		http://127.0.0.1:2379
  advertisePeer:	http://127.0.0.1:2380
  clusterToken:		test-cassini-cluster
  cluster:			test-cassini=http://127.0.0.1:2380
`

	conf := &config.Config{}
	err := conf.Parse([]byte(confStr))
	assert.NoError(t, err)

	viper.Set("useEtcd", true)
	viper.Set("lock", "etcd://127.0.0.1:2379")

	assert.Equal(t, true, viper.GetBool("useEtcd"))

	var etcd *embed.Etcd
	etcd, err = StartEmbedEtcd(conf)
	assert.NoError(t, err)
	defer etcd.Close()

	goroutines := 3

	ms := make([]Mutex, goroutines)

	for i := 0; i < goroutines; i++ {
		m, err := NewMutex("test", viper.GetString("lock"))
		assert.NoError(t, err)
		ms[i] = m
		defer m.Close()
	}

	errs := []error{}
	var sequence int64
	var w sync.WaitGroup
	var mux sync.Mutex

	w.Add(3)
	sequence++

	for i := 0; i < goroutines; i++ {
		c := i
		go func() {
			ms[c].Update(sequence)
			w.Done()
		}()
	}

	w.Wait()

	w.Add(3)

	for i := 0; i < goroutines; i++ {
		c := i
		go func() {
			_, e := ms[c].Lock(sequence)
			if e != nil {
				mux.Lock()
				errs = append(errs, e)
				mux.Unlock()
			} else {
				time.Sleep(1000)
				ms[c].Unlock(true)
			}
			w.Done()
		}()
	}

	w.Wait()
	assert.Equal(t, goroutines-1, len(errs))

	errs = []error{}
	sequence++
	w.Add(3)

	for i := 0; i < goroutines; i++ {
		c := i
		go func() {
			_, e := ms[c].Lock(sequence)
			if e != nil {
				mux.Lock()
				errs = append(errs, e)
				mux.Unlock()
			} else {
				time.Sleep(1000)
				ms[c].Unlock(true)
			}
			w.Done()
		}()
	}

	w.Wait()
	assert.Equal(t, goroutines-1, len(errs))

	os.RemoveAll("./test-cassini.etcd/")

}
