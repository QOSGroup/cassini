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

func main() {

	log.ReplaceConfig(logConfig)

	seq := int64(999)
	err := testEtcdMutexUpdate(seq)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Update sequence: ", seq)

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
	etcd, err = concurrency.NewEtcdMutex("mykey",
		[]string{"192.168.1.195:2379",
			"192.168.1.195:22379",
			"192.168.1.195:32379"})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = etcd.Update(sequence)
	if err != nil {
		fmt.Println("Update error: ", err)
		return
	}
	defer func() {
		err = etcd.Close()
		if err != nil {
			fmt.Println("Close error: ", err)
		}
		fmt.Println("close done")
	}()
	return

}

func testEtcdMutex(sequence int64) (seq int64, err error) {
	seq = -1
	fmt.Printf("Request lock sequence: %d\n", sequence)
	var etcd *concurrency.EtcdMutex
	etcd, err = concurrency.NewEtcdMutex("mykey",
		[]string{"192.168.1.195:2379",
			"192.168.1.195:22379",
			"192.168.1.195:32379"})
	if err != nil {
		fmt.Println(err)
		return
	}
	seq, err = etcd.Lock(sequence)
	defer func() {
		err = etcd.Close()
		if err != nil {
			fmt.Println("Close error: ", err)
		}
		fmt.Println("close done")
	}()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err = etcd.Unlock(true)
		if err != nil {
			fmt.Println("Unlock error: ", err)
		}
		fmt.Println("unlock done")
	}()

	fmt.Printf("Get lock sequence(%d), current sequence(%d) in lock\n", sequence, seq)
	time.Sleep(1 * time.Second)
	return 0, nil

}
