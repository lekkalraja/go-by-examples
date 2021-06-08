* `Go By Examples` : https://gobyexample.com/


### Panic

* A `panic` typically means something went unexpectedly wrong. Mostly we use it to fail fast on errors that shouldn’t occur during normal operation, or that we aren’t prepared to handle gracefully.

* A common use of panic is to abort if a function returns an error value that we don’t know how to (or want to) handle.

* Note that unlike some languages which use exceptions for handling of many errors, in Go it is idiomatic to use error-indicating return values wherever possible.

### Defer

* `Defer` is used to ensure that a function call is performed later in a program’s execution, usually for purposes of cleanup. 
* defer is often used where e.g. `ensure` and `finally` would be used in other languages.

* Suppose we wanted to create a file, write to it, and then close when we’re done. Here’s how we could do that with defer.

    ```go
        func main() {
            f := createFile("/tmp/defer.txt")
            defer closeFile(f)
            writeFile(f)
        }
    ```
* It’s important to check for errors when closing a file, even in a deferred function.

    ```go
        func closeFile(f *os.File) {
            fmt.Println("closing")
            err := f.Close()
            if err != nil {
                fmt.Fprintf(os.Stderr, "error: %v\n", err)
                os.Exit(1)
            }
        }
    ```

### Closures

* Go supports `anonymous functions`, which can form `closures`. 
* Anonymous functions are useful when you want to define a function inline without having to name it.

* The function `intSeq` returns another function, which we define anonymously in the body of intSeq. `The returned function closes over the variable i to form a closure.`

    ```go
        func intSeq() func() int {
            i := 0
            return func() int {
                i++
                return i
            }
        }
    ```
* We call intSeq, assigning the result (a function) to nextInt. This function value captures its own i value, which will be updated each time we call nextInt.

* `To confirm that the state is unique to that particular function, create and test a new one`.

    ```go
        func main() {
            nextInt := intSeq()
            fmt.Println(nextInt())
            fmt.Println(nextInt())
            fmt.Println(nextInt())

            newInts := intSeq()
            fmt.Println(newInts())
        }
    ```

### Errors

* In Go it’s idiomatic to communicate errors via an explicit, separate return value. 
* This contrasts with the exceptions used in languages like Java and Ruby and the overloaded single result / error value sometimes used in C. 
* Go’s approach makes it easy to see which functions return errors and to handle them using the same language constructs employed for any other, non-error tasks.
* By convention, errors are the last return value and have type error, a built-in interface.

* errors.New constructs a basic error value with the given error message.
    ```go
        errors.New("can't work with 42")
    ```
* A nil value in the error position indicates that there was no error.
    ```go
        return arg + 3, nil
    ```

* It’s possible to use custom types as errors by implementing the Error() method on them. 
* Here’s a variant example that uses a custom type to explicitly represent an argument error.

    ```go
        type argError struct {
            arg  int
            prob string
        }

        func (e *argError) Error() string {
            return fmt.Sprintf("%d - %s", e.arg, e.prob)
        }
    ```
* In this case we use &argError syntax to build a new struct, supplying values for the two fields arg and prob.
    ```go
        return -1, &argError{arg, "can't work with it"}
    ```

* The loop below test out each of our error-returning functions. Note that the use of an inline error check on the if line is a common idiom in Go code.

    ```go
        for _, i := range []int{7, 42} {
            if r, e := f1(i); e != nil {
                fmt.Println("f1 failed:", e)
            } else {
                fmt.Println("f1 worked:", r)
            }
        }
    ```
* If you want to programmatically use the data in a custom error, you’ll need to get the error as an instance of the custom error type `via type assertion`.

    ```go
        _, e := f2(42)
        if ae, ok := e.(*argError); ok {
            fmt.Println(ae.arg)
            fmt.Println(ae.prob)
        }
    ```
* Great Blog on Error Handling : https://blog.golang.org/error-handling-and-go

### Constants


* Go supports constants of character, string, boolean, and numeric values.

* `const` declares a constant value.
* A const statement can appear anywhere a var statement can.
* Constant expressions perform arithmetic with arbitrary precision.
* A numeric constant has no type until it’s given one, such as by an explicit conversion.
* A number can be given a type by using it in a context that requires one, such as a variable assignment or function call. For example, here math.Sin expects a float64.

    ```go
        const s string = "constant"

        const n = 500000000

        const d = 3e20 / n

        fmt.Println(int64(d))

        fmt.Println(math.Sin(n))
    ```

### HTTP Servers

* Writing a basic HTTP server is easy using the `net/http` package.

* A fundamental concept in net/http servers is `handlers`. 
* A handler is an object implementing the `http.Handler interface`.
*  A common way to write a handler is by using the `http.HandlerFunc` adapter on functions with the appropriate signature.

* Functions serving as handlers take a http.ResponseWriter and a http.Request as arguments.
* The response writer is used to fill in the HTTP response. 

    ```go
        func hello(w http.ResponseWriter, req *http.Request) {
            fmt.Fprintf(w, "hello\n")
        }
    ```

* We register our handlers on server routes using the http.HandleFunc convenience function. It sets up the default router in the `net/http` package and takes a function as an argument.

* Finally, we call the `ListenAndServe` with the port and a handler. `nil tells it to use the default router we’ve just set up`.

    ```go
        http.HandleFunc("/hello", hello)
        http.ListenAndServe(":8090", nil)
    ```

### HTTP Clients

* The Go standard library comes with excellent support for HTTP clients and servers in the `net/http` package.

* `http.Get` is a convenient shortcut around creating an `http.Client` object and calling its Get method. `it uses the `http.DefaultClient` object which has useful default settings.

    ```go
        resp, err := http.Get("http://gobyexample.com")
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
    ```

### Context

* we have seen a simple HTTP server. HTTP servers are useful for demonstrating the usage of `context.Context` for controlling cancellation. 
* A Context carries deadlines, cancellation signals, and other request-scoped values across API boundaries and goroutines.

* A `context.Context` is created for each request by the `net/http` machinery, and is available with the `Context()` method.

* Wait for a few seconds before sending a reply to the client. This could simulate some work the server is doing. 
* While working, keep an eye on the `context’s Done()` channel for a signal that we should cancel the work and return as soon as possible.
* The context’s Err() method returns an error that explains why the Done() channel was closed.

    ```go
        func hello(w http.ResponseWriter, req *http.Request) {
            ctx := req.Context()
            fmt.Println("server: hello handler started")
            defer fmt.Println("server: hello handler ended")
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

    ```
* As before, we register our handler on the “/hello” route, and start serving.

    ```go
        http.HandleFunc("/hello", hello)
        http.ListenAndServe(":8090", nil)
    ```
### Spawning Processes

* Sometimes our Go programs need to spawn other, `non-Go processes`.

* a simple command that takes no arguments or input and just prints something to stdout. The    `exec.Command` helper creates an object to represent this external process.
* `.Output` is another helper that handles the common case of running a command, waiting for it to finish, and collecting its output. If there were no errors, dateOut will hold bytes with the date info.

    ```go
        dateCmd := exec.Command("date")

        dateOut, err := dateCmd.Output()
        if err != nil {
            panic(err)
        }
        fmt.Println("> date")
        fmt.Println(string(dateOut))
    ```
* Next we’ll look at a slightly more involved case where we pipe data to the external process on its stdin and collect the results from its stdout.

* Here we explicitly grab input/output pipes, start the process, write some input to it, read the resulting output, and finally wait for the process to exit.

    ```go
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
    ```

* Note that when spawning commands we need to provide an explicitly delineated command and argument array, vs. being able to just pass in one command-line string.
* If you want to spawn a full command with a string, you can use bash’s -c option:
* The spawned programs return output that is the same as if we had run them directly from the command-line.

    ```go
        lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
        lsOut, err := lsCmd.Output()
        if err != nil {
            panic(err)
        }
        fmt.Println("> ls -a -l -h")
        fmt.Println(string(lsOut))
    ```

### Exec'ing Processes

* Earlier we looked at spawning external processes. We do this when we need an external process accessible to a running Go process. 
* Sometimes we just want to completely replace the current Go process with another (perhaps non-Go) one. 
* To do this we’ll use Go’s implementation of the `classic exec function`.

* For our example we’ll exec ls. Go requires an absolute path to the binary we want to execute, so we’ll use exec.LookPath to find it (probably /bin/ls).

    ```go
        binary, lookErr := exec.LookPath("ls")
        if lookErr != nil {
            panic(lookErr)
        }
    ```
* Exec requires arguments in slice form (as opposed to one big string). We’ll give ls a few common arguments. `Note that the first argument should be the program name`.
* Exec also needs a set of environment variables to use. Here we just provide our current environment.
* Here’s the actual syscall.Exec call. `If this call is successful, the execution of our process will end here and be replaced by the /bin/ls -a -l -h process`. 
* If there is an error we’ll get a return value.

    ```go
        args := []string{"ls", "-a", "-l", "-h"}
        env := os.Environ()
        execErr := syscall.Exec(binary, args, env)
        if execErr != nil {
            panic(execErr)
        }
    ```

* When we run our program it is replaced by ls.

* Note that Go does not offer a classic Unix fork function. Usually this isn’t an issue though, since starting goroutines, spawning processes, and exec’ing processes covers most use cases for fork.

### Signals
* Sometimes we’d like our Go programs to intelligently `handle Unix signals`. 
* For example, we might want a server to `gracefully shutdown` when it receives a `SIGTERM`, or a command-line tool to stop processing input if it receives a `SIGINT`. 
* Here’s how to handle signals in Go with channels.
* Go signal notification works by sending `os.Signal` values on a channel. We’ll create a channel to receive these notifications (we’ll also make one to notify us when the program can exit).

    ```go
        sigs := make(chan os.Signal, 1)
        done := make(chan bool, 1)
    ```
* signal.Notify registers the given channel to receive notifications of the specified signals.

    ```go
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    ```

* This goroutine executes a blocking receive for signals. When it gets one it’ll print it out and then notify the program that it can finish.

    ```go
        go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println(sig)
        done <- true
    }()
    ```
* The program will wait here until it gets the expected signal (as indicated by the goroutine above sending a value on done) and then exit.

    ```go
        fmt.Println("awaiting signal")
        <-done
        fmt.Println("exiting")
    ```
* When we run this program it will block waiting for a signal. By typing ctrl-C (which the terminal shows as ^C) we can send a SIGINT signal, causing the program to print interrupt and then exit.

### Exit

* Use `os.Exit` to immediately exit with a given status.

* `defers will not be run when using os.Exit`, so this fmt.Println will never be called.
* Exit with status 3.
    ```go
        defer fmt.Println("defer !")
        os.Exit(3)
    ```
* Note that unlike e.g. C, Go does not use an integer return value from main to indicate exit status. If you’d like to exit with a non-zero status you should use os.Exit.

* If you run exit.go using go run, the exit will be picked up by go and printed.
    ```go
        $ go run exit.go
        exit status 3
    ```

* By building and executing a binary you can see the status in the terminal.

    ```go
        $ go build exit.go
        $ ./exit
        $ echo $?
        3
    ```
* Note that the `defer !` from our program `never got printed`.