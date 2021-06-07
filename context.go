package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	log.Println("server: hello handler started")
	defer log.Println("server: hello handler ended")

	select {
	case <-time.After(10 * time.Second):
		fmt.Fprintf(w, "hello\n")
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("server:", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8090", nil)
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run context.go
	2021/06/07 10:26:37 server: hello handler started
	2021/06/07 10:26:47 server: hello handler ended
	2021/06/07 10:26:51 server: hello handler started
	server: context canceled (when we cancel the client call)
	2021/06/07 10:26:53 server: hello handler ended
	^Csignal: interrupt
*/
