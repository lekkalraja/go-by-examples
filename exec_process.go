package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {

	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"ls", "-a", "-l", "-h"}

	env := os.Environ()

	//log.Printf("Environment : %s\n", env)

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run exec_process.go
	total 76K
	drwxrwxr-x 3 raja raja 4.0K Jun  7 10:53 .
	drwxrwxr-x 7 raja raja 4.0K Jun  6 10:05 ..
	-rw-rw-r-- 1 raja raja  747 Jun  6 10:58 closures.go
	-rw-rw-r-- 1 raja raja  359 Jun  6 11:36 constants.go
	-rw-rw-r-- 1 raja raja  938 Jun  7 10:27 context.go
	-rw-rw-r-- 1 raja raja  905 Jun  6 10:46 defer.go
	-rw-rw-r-- 1 raja raja 1.1K Jun  6 11:11 errors.go
	-rw-rw-r-- 1 raja raja  345 Jun  7 10:54 exec_process.go
	drwxrwxr-x 8 raja raja 4.0K Jun  7 10:44 .git
	-rw-rw-r-- 1 raja raja  269 Jun  6 10:05 .gitignore
	-rw-rw-r-- 1 raja raja 9.7K Jun  7 09:46 http_client.go
	-rw-rw-r-- 1 raja raja 1.1K Jun  7 09:38 http_server.go
	-rw-rw-r-- 1 raja raja 1.4K Jun  6 10:38 panic.go
	-rw-rw-r-- 1 raja raja  12K Jun  7 10:53 README.md
	-rw-rw-r-- 1 raja raja 1.6K Jun  7 10:44 spawing_process.go
*/
