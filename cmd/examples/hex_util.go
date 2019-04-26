// +build ignore

package main

import (
	"fmt"
	"math/big"
	"strconv"
)

func main() {

	gas := "0x3777c02e70512800" //"0x3782dace9d900000"

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

	gas = "0xc980"
	gasPrice = "0x12a05f200"
	fee := count(gas, gasPrice)

	gas = "0x8ee8"
	gasPrice = "0x12a05f200"
	tmp := count(gas, gasPrice)
	fee.Add(fee, tmp)
	fmt.Println("fee: ", fee.Text(16), "; ", fee.Text(10))

	gas = "0x90a2"
	gasPrice = "0x2540be400"
	tmp = count(gas, gasPrice)
	fee.Add(fee, tmp)
	fmt.Println("fee: ", fee.Text(16), "; ", fee.Text(10))

	gas = "0x90a2"
	gasPrice = "0x2540be400"
	tmp = count(gas, gasPrice)
	fee.Add(fee, tmp)
	fmt.Println("fee: ", fee.Text(16), "; ", fee.Text(10))

	gas = "0x10ae2"
	gasPrice = "0x2540be400"
	tmp = count(gas, gasPrice)
	fee.Add(fee, tmp)
	fmt.Println("fee: ", fee.Text(16), "; ", fee.Text(10))

	balance := sum("0xde0b6b3a7640000", "0x29a2241af62c0000")
	fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))
	balance.Sub(balance, fee)
	fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))

	sub(balance, "0xbf69898a1800")
	fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))

	sub(balance, "0xbf69898a1800")
	fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))

	sub(balance, "0x17dfcdece4000")
	fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))

	sub(balance, "0xbefe6f672000")
	fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))

	sub(balance, "0xbefe6f672000")
	fmt.Println("last balance: ", balance.Text(16), "; ", balance.Text(10))

	sub(balance, "0xbefe6f672000")
	fmt.Println("last balance: ", balance.Text(16), "; ", balance.Text(10))
	// ////////////////////

	// gas = "0x12e40"
	// gasPrice = "0x12a05f200"
	// f := count(gas, gasPrice)

	// gas = "0xd65c"
	// gasPrice = "0x12a05f200"
	// tmp = count(gas, gasPrice)
	// f.Add(f, tmp)
	// fmt.Println("=== fee: ", f.Text(16), "; ", f.Text(10))

	// gas = "0xd8f3"
	// gasPrice = "0x2540be400"
	// tmp = count(gas, gasPrice)
	// f.Add(f, tmp)
	// fmt.Println("fee: ", f.Text(16), "; ", f.Text(10))

	// gas = "0xd8f3"
	// gasPrice = "0x2540be400"
	// tmp = count(gas, gasPrice)
	// f.Add(f, tmp)
	// fmt.Println("fee: ", f.Text(16), "; ", f.Text(10))

	// gas = "0x19053"
	// gasPrice = "0x2540be400"
	// tmp = count(gas, gasPrice)
	// f.Add(f, tmp)
	// fmt.Println("fee: ", f.Text(16), "; ", f.Text(10))

	// balance = sum("0xde0b6b3a7640000", "0x29a2241af62c0000")
	// fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))
	// balance.Sub(balance, f)
	// fmt.Println("balance: ", balance.Text(16), "; ", balance.Text(10))

	balance = sum("0x00f", "0x-3")
	fmt.Println("sum balance: ", balance.Text(16), "; ", balance.Text(10))
}

func count(gasUsed, gasPrice string) *big.Int {
	g := new(big.Int)
	g.SetString(gasUsed[2:], 16)
	gp := new(big.Int)
	gp.SetString(gasPrice[2:], 16)
	g = g.Mul(g, gp)
	// fee := fmt.Sprintf("0x%s\n", g.Text(16))
	// fmt.Println("fee: ", fee, "; ", g.Text(10))
	return g
}

func sum(gasUsed, gasPrice string) *big.Int {
	g := new(big.Int)
	g.SetString(gasUsed[2:], 16)
	gp := new(big.Int)
	gp.SetString(gasPrice[2:], 16)
	g.Add(g, gp)
	return g
}

func sub(i *big.Int, v string) *big.Int {
	g := new(big.Int)
	g.SetString(v[2:], 16)
	i.Sub(i, g)
	return i
}
