package main

import (
	"context"
	"fmt"

	"github.com/QOSGroup/cassini/config"
)

// command reset handler
var resetHandler = func(conf *config.Config) (cancel context.CancelFunc, err error) {
	fmt.Println("!!!WARN It's DANGER!!! Data will be reset(cleaning up)!")

	return
}
