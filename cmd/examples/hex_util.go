// +build ignore

package main

import (
	"fmt"
	"strconv"
)

func main() {

	s := "0x3782dace9d900000"

	v, err := strconv.ParseInt(s[2:], 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s, " int ", v)
}
