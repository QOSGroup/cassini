package common

import (
	"github.com/QOSGroup/cassini/log"
	"os"
	"os/signal"
	"syscall"
)

// KeepRunning 保持程序运行，监听系统信号，触发回调函数
func KeepRunning(callback func(sig os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case s, ok := <-signals:
		log.Infof("System signal [%v] %t, trying to run callback...", s, ok)
		if !ok {
			break
		}
		if callback != nil {
			callback(s)
		}
		os.Exit(1)
	}
}
