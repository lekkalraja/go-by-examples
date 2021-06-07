package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {

	dateCmd := exec.Command("date")

	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))

	grepCmd := exec.Command("grep", "hello")

	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	grepIn.Close()
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()

	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}

/*
raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run spawing_process.go
> date
Mon  7 Jun 10:43:20 +08 2021

> grep hello
hello grep

> ls -a -l -h
total 72K
drwxrwxr-x 3 raja raja 4.0K Jun  7 10:42 .
drwxrwxr-x 7 raja raja 4.0K Jun  6 10:05 ..
-rw-rw-r-- 1 raja raja  747 Jun  6 10:58 closures.go
-rw-rw-r-- 1 raja raja  359 Jun  6 11:36 constants.go
-rw-rw-r-- 1 raja raja  938 Jun  7 10:27 context.go
-rw-rw-r-- 1 raja raja  905 Jun  6 10:46 defer.go
-rw-rw-r-- 1 raja raja 1.1K Jun  6 11:11 errors.go
drwxrwxr-x 8 raja raja 4.0K Jun  7 10:28 .git
-rw-rw-r-- 1 raja raja  269 Jun  6 10:05 .gitignore
-rw-rw-r-- 1 raja raja 9.7K Jun  7 09:46 http_client.go
-rw-rw-r-- 1 raja raja 1.1K Jun  7 09:38 http_server.go
-rw-rw-r-- 1 raja raja 1.4K Jun  6 10:38 panic.go
-rw-rw-r-- 1 raja raja  11K Jun  7 10:42 README.md
-rw-rw-r-- 1 raja raja  718 Jun  7 10:42 spawing_process.go
*/
