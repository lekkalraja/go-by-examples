package main

import (
	"fmt"
	"math"
)

const s string = "constant"

func main() {
	fmt.Println(s)

	const n = 500000000

	const d = 3e20 / n
	fmt.Println(d)

	fmt.Println(int64(d))

	fmt.Println(math.Sin(n))
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run constants.go
	constant
	6e+11
	600000000000
	-0.28470407323754404
*/
