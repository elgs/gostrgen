# gostrgen
Random string generator in Golang.

#Installation
`go get -u github.com/elgs/gostrgen`

# Sample code
```go
package main

import (
	"fmt"
	"github.com/elgs/gostrgen"
)

func main() {

	// possible alphabet set are:
	// Lower, Upper, Digit, Punct, LowerUpper, LowerDigit, UpperDigit, LowerUpperDigit and All.
	// Any of the above can be combine by "|", e.g. LowerUpper is the same as Lower | Upper

	chararactersToGenerate := 20
	set := gostrgen.Lower | gostrgen.Digit
	include := "[]{}<>" // optionally include some additional letters
	exclude := "Ol"     //exclude big 'O' and small 'l' to avoid confusion with zero and one.

	str, err := gostrgen.RandGen(chararactersToGenerate, set, include, exclude)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str) // zxh9[pvoxbaup32b7s0d
}
```