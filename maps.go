package main

import "fmt"

func main() {

	m := make(map[string]int)

	m["k1"] = 7
	m["k2"] = 13

	fmt.Println("map:", m)

	v1 := m["k1"]
	fmt.Println("v1: ", v1)

	fmt.Println("len:", len(m))

	delete(m, "k2")
	fmt.Println("map:", m)

	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run maps.go
	map: map[k1:7 k2:13]
	v1:  7
	len: 2
	map: map[k1:7]
	prs: false
	map: map[bar:2 foo:1]
*/
