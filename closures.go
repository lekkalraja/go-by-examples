package main

import (
	"log"
)

func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {

	nextInt := intSeq()

	log.Printf("[1] first increment of i value : %d", nextInt())
	log.Printf("[1] second increment of i value : %d", nextInt())
	log.Printf("[1] third increment of i value : %d", nextInt())

	anotherInt := intSeq()

	log.Printf("[2] first increment of i value : %d", anotherInt())
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run closures.go
	2021/06/06 10:57:58 [1] first increment of i value : 1
	2021/06/06 10:57:58 [1] second increment of i value : 2
	2021/06/06 10:57:58 [1] third increment of i value : 3
	2021/06/06 10:57:58 [2] first increment of i value : 1
*/
