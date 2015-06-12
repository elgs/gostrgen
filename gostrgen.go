// gostrgen
package gostrgen

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

var Lower = 1 << 0
var Upper = 1 << 1
var Digit = 1 << 2
var Punct = 1 << 3

var LowerUpper = Lower | Upper
var LowerDigit = Lower | Digit
var UpperDigit = Upper | Digit
var LowerUpperDigit = LowerUpper | Digit
var All = LowerUpperDigit | Punct

var lower = "abcdefghijklmnopqrstuvwxyz"
var upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var digit = "0123456789"
var punct = "~!@#$%^&*()_+-="

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandGen(size int, set int, include string, exclude string) (string, error) {
	all := include
	if set&Lower > 0 {
		all += lower
	}
	if set&Upper > 0 {
		all += upper
	}
	if set&Digit > 0 {
		all += digit
	}
	if set&Punct > 0 {
		all += punct
	}

	lenAll := len(all)
	if len(exclude) >= lenAll {
		return "", errors.New("Too much to exclude.")
	}
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		b := all[rand.Intn(lenAll)]
		if strings.Contains(exclude, string(b)) {
			i--
			continue
		}
		buf[i] = b
	}
	return string(buf), nil
}
