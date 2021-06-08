package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run signals.go
	awaiting signal
	^C
	interrupt
	exiting
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$
*/
