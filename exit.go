package main

import (
	"fmt"
	"os"
)

func main() {

	defer fmt.Println("defer !")

	os.Exit(3)
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run exit.go
	exit status 3
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go build exit.go
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ ./exit
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ echo $?
	3
*/
