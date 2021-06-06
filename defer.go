package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f := createFile("/tmp/defer.txt")
	defer closeFile(f)
	writeFile(f)
}

func createFile(p string) *os.File {
	log.Printf("creating file : %s \n", p)
	f, err := os.Create(p)
	if err != nil {
		log.Panicf("Something went wrong while creating file : %v \n", err)
	}
	return f
}

func writeFile(f *os.File) {
	log.Printf("writing data to file : %s \n", f.Name())
	fmt.Fprintln(f, "data")
}

func closeFile(f *os.File) {
	log.Printf("closing file : %s \n", f.Name())
	err := f.Close()

	if err != nil {
		log.Fatalf("something went wrong while closing file : %s , error: %v\n", f.Name(), err)
	}
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run defer.go
	2021/06/06 10:45:57 creating file : /tmp/defer.txt
	2021/06/06 10:45:57 writing data to file : /tmp/defer.txt
	2021/06/06 10:45:57 closing file : /tmp/defer.txt
*/
