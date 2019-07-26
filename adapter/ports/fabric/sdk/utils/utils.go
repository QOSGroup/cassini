/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0

https://github.com/securekey/fabric-examples/blob/master/fabric-cli/cmd/fabric-cli/chaincode/utils/util.go

*/

package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

const (
	randFunc = "$rand("
	padFunc  = "$pad("
	seqFunc  = "$seq("
	setFunc  = "$set("
	varExp   = "${"
)

var (
	sequence uint64
)

type Context interface {
	SetVar(name, value string)
	GetVar(name string) (string, bool)
}

// AsBytes converts the string array to an array of byte arrays.
// The args may contain functions $rand(n) or $pad(n,chars).
// The functions are evaluated before returning.
//
// Examples:
// - "key$rand(3)" -> "key0" or "key1" or "key2"
// - "val$pad(3,XYZ) -> "valXYZXYZXYZ"
// - "val$pad($rand(3),XYZ) -> "val" or "valXYZ" or "valXYZXYZ"
// - "key$seq()" -> "key1", "key2", "key2", ...
// - "val$pad($seq(),X)" -> "valX", "valXX", "valXX", "valXXX", ...
// - "Key_$set(x,$seq())=Val_${x}" -> Key_1=Val_1, Key_2=Val_2, ...
func AsBytes(args []string) [][]byte {
	rand := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	bytes := make([][]byte, len(args))

	// TODO config verbose?
	verbose := false

	if verbose {
		fmt.Printf("Args:\n")
	}
	for i, a := range args {
		arg := getArg(rand, a)
		if verbose {
			fmt.Printf("- [%d]=%s\n", i, arg)
		}
		bytes[i] = []byte(arg)
	}
	return bytes
}

func getArg(r *rand.Rand, arg string) string {
	arg = evaluateSeqExpression(arg)
	arg = evaluateRandExpression(r, arg)
	arg = evaluatePadExpression(arg)
	return arg
}

// evaluateSeqExpression replaces occurrences of $seq() with a sequential
// number starting at 1 and incrementing for each task
func evaluateSeqExpression(arg string) string {
	return evaluateExpression(arg, seqFunc, ")",
		func(expression string) (string, error) {
			return strconv.FormatUint(atomic.AddUint64(&sequence, 1), 10), nil
		})
}

// evaluateRandExpression replaces occurrences of $rand(n) with a random
// number between 0 and n (exclusive)
func evaluateRandExpression(r *rand.Rand, arg string) string {
	return evaluateExpression(arg, randFunc, ")",
		func(expression string) (string, error) {
			n, err := strconv.ParseInt(expression, 10, 64)
			if err != nil {
				return "", errors.Errorf("invalid number %s in $rand expression\n", expression)
			}
			return strconv.FormatInt(r.Int63n(n), 10), nil
		})
}

// evaluatePadExpression replaces occurrences of $pad(n,chars) with n of the given pad characters
func evaluatePadExpression(arg string) string {
	return evaluateExpression(arg, padFunc, ")",
		func(expression string) (string, error) {
			s := strings.Split(expression, ",")
			if len(s) != 2 {
				return "", errors.Errorf("invalid $pad expression: '%s'. Expecting $pad(n,chars)", expression)
			}

			n, err := strconv.Atoi(s[0])
			if err != nil {
				return "", errors.Errorf("invalid number %s in $pad expression\n", s[0])
			}

			result := ""
			for i := 0; i < n; i++ {
				result += s[1]
			}

			return result, nil
		})
}

func evaluateExpression(expression, funcType, endDelim string, evaluate func(string) (string, error)) string {
	result := ""
	for {
		i := strings.Index(expression, funcType)
		if i == -1 {
			result += expression
			break
		}

		j := strings.Index(expression[i:], endDelim)
		if j == -1 {
			fmt.Printf("expecting '%s' in expression '%s'", endDelim, expression)
			result = expression
			break
		}

		j = i + j

		replacement, err := evaluate(expression[i+len(funcType) : j])
		if err != nil {
			fmt.Printf("%s\n", err)
			result += expression[0 : j+1]
		} else {
			result += expression[0:i] + replacement
		}

		expression = expression[j+1:]
	}

	return result
}
