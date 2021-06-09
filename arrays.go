package main

import "fmt"

func main() {

	var a [5]int
	fmt.Println("emp:", a)

	a[4] = 100
	fmt.Println("set:", a)
	fmt.Println("get:", a[4])

	fmt.Println("len:", len(a))

	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("dcl:", b)

	update(b) // Will not change because it is pass-by-value

	fmt.Printf("[ %p ]After Updting: %v\n", &b, b)

	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}

func update(arr [5]int) {
	fmt.Printf("[ %p ]Updting: %v\n", &arr, arr)
	arr[3] = 132
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run arrays.go
	emp: [0 0 0 0 0]
	set: [0 0 0 0 100]
	get: 100
	len: 5
	dcl: [1 2 3 4 5]
	[ 0xc0001800f0 ]Updting: [1 2 3 4 5]
	[ 0xc000180090 ]After Updting: [1 2 3 4 5]
	2d:  [[0 1 2] [1 2 3]]
*/
