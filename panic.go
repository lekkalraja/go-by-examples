package main

import (
	"log"
	"os"
)

func main() {
	_, err := os.Create("/tmp/abc/xyz/file.txt")

	if err != nil {
		// panic(err)
		// log.Panicf("Something went wrong while creating file , %v", err)
		log.Fatalf("Something went wrong while creating file , %v", err)
	}
}

/*

	* Running this program will cause it to panic,
	  print an error message and goroutine traces, and exit with a non-zero status.

	panic: open /tmp/abc/xyz/file.txt: no such file or directory

	goroutine 1 [running]:
	main.main()
        /home/raja/Documents/coding/golang/go-by-examples/panic.go:9 +0x7a
	exit status 2

	log.Panicf : equivalent to Printf() and then panic()
	====================================================

		2021/06/06 10:35:53 Something went wrong while creating file , open /tmp/abc/xyz/file.txt: no such file or directory
	panic: Something went wrong while creating file , open /tmp/abc/xyz/file.txt: no such file or directory

	goroutine 1 [running]:
	log.Panicf(0x4ccfb1, 0x2d, 0xc000068f68, 0x1, 0x1)
			/usr/local/go/src/log/log.go:361 +0xc5
	main.main()
			/home/raja/Documents/coding/golang/go-by-examples/panic.go:13 +0xad
	exit status 2

	log.Fatalf: equivalent to Printf() and then os.Exit(1)
	=====================================================

	2021/06/06 10:37:15 Something went wrong while creating file , open /tmp/abc/xyz/file.txt: no such file or directory
	exit status 1
*/
