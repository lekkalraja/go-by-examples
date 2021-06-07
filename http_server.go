package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Panicf("Something went wrong while starting Http Server : %v \n", err)
	}

}

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello, World!\n")
}

func headers(resp http.ResponseWriter, req *http.Request) {
	for key, value := range req.Header {
		fmt.Fprintf(resp, "Header key: %s, value : %s\n", key, value)
	}
}

/*
	Run The Server : go run http_server.go


	Use CURL to hit the endpoints:
	==============================
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ curl http://localhost:8080/hello
	Hello, World!
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ curl http://localhost:8080/headers
	Header key: User-Agent, value : [curl/7.68.0]
	Header key: Accept, value : [*/ /*] (modified form / to // to escape the comment)
raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$

*/
