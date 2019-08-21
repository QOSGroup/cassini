package main

import (
	"context"
	"fmt"
)

// command reset handler
var resetHandler = func() (cancel context.CancelFunc, err error) {
	fmt.Println("!!!WARN It's DANGER!!! Data will be reset(cleaning up)!")

	return
}
