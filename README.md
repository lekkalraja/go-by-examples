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