package main

import (
	"fmt"
	"os"
	"time"

	"github.com/QOSGroup/cassini/concurrency"
	"github.com/QOSGroup/cassini/log"

	cmn "github.com/QOSGroup/cassini/common"
)

var logConfig = `
<seelog minlevel="off">
	<outputs formatid="formater"><console />
	</outputs>
	<formats>
		<format id="formater" format="[%Date(2006-01-02 15:04:05.000000000)][%LEV] %Msg%n"/>
	</formats>
</seelog>`

var mykey = "mykey"

// var addrs = []string{"192.168.1.195:2379", "192.168.1.195:22379", "192.168.1.195:32379"}
var addrs = []string{"192.168.1.35:2379"}

func main() {

	log.ReplaceConfig(logConfig)

	seq := int64(999)
	err := testEtcdMutexUpdate(seq)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Update sequence: ", seq)

	testEtcdMutexLoop(3)
	fmt.Println("Started loop test")

	go func() {
		var c int64
		for {
			for i := 0; i < 3; i++ {
				c++
				s := c
				go func() {
					seq, _ := testEtcdMutex(s)
					if seq > -1 {
						c = seq - 1
					}
				}()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	cmn.KeepRunning(func(sig os.Signal) {
		fmt.Println("OK, bye.")
	})
}

func testEtcdMutexUpdate(sequence int64) (err error) {
	fmt.Printf("Request lock sequence: %d\n", sequence)
	var etcd *concurrency.EtcdMutex
	etcd, err = concurrency.NewEtcdMutex(mykey, addrs)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = etcd.Update(sequence)
	if err != nil {
		log.Warn("Update error: ", err)
		return
	}
	defer func() {
		err = etcd.Close()
		if err != nil {
			log.Warn("Close error: ", err)
		}
		log.Debug("close done")
	}()
	return

}

func testEtcdMutex(sequence int64) (seq int64, err error) {
	seq = -1
	log.Debugf("Request lock sequence: %d\n", sequence)
	var etcd *concurrency.EtcdMutex
	etcd, err = concurrency.NewEtcdMutex(mykey, addrs)
	if err != nil {
		fmt.Println(err)
		return
	}
	seq, err = etcd.Lock(sequence)
	defer func() {
		err = etcd.Close()
		if err != nil {
			log.Warn("Close error: ", err)
		}
		log.Debug("close done")
	}()
	if err != nil {
		log.Warn(err)
		return
	}
	defer func() {
		err = etcd.Unlock(true)
		if err != nil {
			log.Warn("Unlock error: ", err)
		}
		log.Debug("unlock done")
	}()

	fmt.Printf("Get lock sequence(%d), current sequence(%d) in lock\n", sequence, seq)
	time.Sleep(1 * time.Second)
	return 0, nil

}

func testEtcdMutexLoop(c int) {

	for i := 0; i < c; i++ {
		go func() {
			etcd, err := concurrency.NewEtcdMutex(mykey, addrs)
			if err != nil {
				log.Warn(err)
				return
			}

			var sequence int64
			for {
				log.Debugf("Loop request lock sequence: %d\n", sequence)
				sequence, err = etcd.Lock(sequence)
				if err != nil {
					log.Warn("Loop error: ", err)
					continue
				}

				fmt.Printf("Loop get lock sequence(%d)!!!!!!!!!!!!!!!!\n", sequence)
				time.Sleep(1 * time.Second)
				err = etcd.Unlock(true)
				if err != nil {
					log.Warn("Loop unlock error: ", err)
				}
				log.Debug("Loop unlock done")
			}
			// err = etcd.Close()
			// if err != nil {
			// 	fmt.Println("Close error: ", err)
			// }
			// fmt.Println("Loop close done")
		}()
	}

}
