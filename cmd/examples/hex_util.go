// +build ignore

package main

import (
	"fmt"
	"strconv"
)

func main() {

	gas := "0x12e40" //"0x3782dace9d900000"

	gasV, err := strconv.ParseInt(gas[2:], 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("gas: ", gas, " int ", gasV)

	gasPrice := "0x12a05f200" //"0x3782dace9d900000"

	gasPriceV, err := strconv.ParseInt(gasPrice[2:], 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("gasPrice: ", gasPrice, " int ", gasPriceV)

	fmt.Println(gasPriceV * gasV)

	s := "0x3777c02e70512800"
	v, err := strconv.ParseInt(s[2:], 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s, " int ", v)
}
