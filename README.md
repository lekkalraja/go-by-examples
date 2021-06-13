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

* Here the variable `i` `Escapes to the Heap`

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

### JSON

* Go offers built-in support for JSON encoding and decoding, including to and from built-in and custom data types.

* We’ll use these two structs to demonstrate encoding and decoding of custom types below.
* Only exported fields will be encoded/decoded in JSON. Fields must start with capital letters to be exported.

    ```go
        type response1 struct {
            Page   int
            Fruits []string
        }

        type response2 struct {
            Page   int      `json:"page"`
            Fruits []string `json:"fruits"`
        }
    ```

* First we’ll look at encoding basic data types to JSON strings. Here are some examples for atomic values.

    ```go
        bolB, _ := json.Marshal(true)
        fmt.Println(string(bolB))
        intB, _ := json.Marshal(1)
        fmt.Println(string(intB))
        fltB, _ := json.Marshal(2.34)
        fmt.Println(string(fltB))
        strB, _ := json.Marshal("gopher")
        fmt.Println(string(strB))
    ```
* here are some for slices and maps, which encode to JSON arrays and objects as you’d expect.

    ```go
        slcD := []string{"apple", "peach", "pear"}
        slcB, _ := json.Marshal(slcD)
        fmt.Println(string(slcB))
        mapD := map[string]int{"apple": 5, "lettuce": 7}
        mapB, _ := json.Marshal(mapD)
        fmt.Println(string(mapB))
    ```

* The JSON package can automatically encode your custom data types. 
* It will only include exported fields in the encoded output and will `by default use those names as the JSON keys`.

    ```go
        res1D := &response1{
            Page:   1,
            Fruits: []string{"apple", "peach", "pear"}
        }
        res1B, _ := json.Marshal(res1D)
        fmt.Println(string(res1B))
    ```

* You can use tags on struct field declarations to customize the encoded JSON key names. Check the definition of response2 above to see an example of such tags.

    ```go
        res2D := &response2{
            Page:   1,
            Fruits: []string{"apple", "peach", "pear"}
        }
        res2B, _ := json.Marshal(res2D)
        fmt.Println(string(res2B))
    ```

* Now let’s look at decoding JSON data into Go values. Here’s an example for a generic data structure.

    ```go
        byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
    ```

* We need to provide a variable where the JSON package can put the decoded data. 
* This map[string]interface{} will hold a map of strings to arbitrary data types.

    ```go
        var dat map[string]interface{}

        if err := json.Unmarshal(byt, &dat); err != nil {
            panic(err)
        }

        fmt.Println(dat)
    ```
* In order to use the values in the decoded map, we’ll need to convert them to their appropriate type.
* For example here we convert the value in num to the expected float64 type.

    ```go
        num := dat["num"].(float64)
        fmt.Println(num)
    ```
* Accessing nested data requires a series of conversions.

    ```go
        strs := dat["strs"].([]interface{})
        str1 := strs[0].(string)
        fmt.Println(str1)
    ```
* We can also decode JSON into custom data types. This has the advantages of adding additional type-safety to our programs and eliminating the need for type assertions when accessing the decoded data.

    ```go
        str := `{"page": 1, "fruits": ["apple", "peach"]}`
        res := response2{}
        json.Unmarshal([]byte(str), &res)
        fmt.Println(res)
        fmt.Println(res.Fruits[0])
    ```

* In the examples above we always used bytes and strings as intermediates between the data and JSON representation on standard out.
* We can also stream JSON encodings directly to os.Writers like os.Stdout or even HTTP response bodies.

    ```go
        enc := json.NewEncoder(os.Stdout)
        d := map[string]int{"apple": 5, "lettuce": 7}
        enc.Encode(d)
    ```

### XML
* Go offers built-in support for XML and XML-like formats with the encoding.xml package.

* Plant will be mapped to XML. Similarly to the JSON examples, field tags contain directives for the encoder and decoder.
* Here we use some special features of the XML package: the `XMLName` field name dictates the name of the XML element representing this struct; `id,attr` means that the Id field is an XML attribute rather than a nested element.

    ```go
        type Plant struct {
            XMLName xml.Name `xml:"plant"`
            Id      int      `xml:"id,attr"`
            Name    string   `xml:"name"`
            Origin  []string `xml:"origin"`
        }

        func (p Plant) String() string {
            return fmt.Sprintf("Plant id=%v, name=%v, origin=%v",p.Id, p.Name, p.Origin)
        }
    ```

* Emit XML representing our plant; using MarshalIndent to produce a more human-readable output.

    ```go
        coffee := &Plant{Id: 27, Name: "Coffee"}
        coffee.Origin = []string{"Ethiopia", "Brazil"}


        out, _ := xml.MarshalIndent(coffee, " ", "  ")
        fmt.Println(string(out))
    ```

* To add a generic XML header to the output, append it explicitly.

    ```go
        fmt.Println(xml.Header + string(out))
    ```

* Use Unmarhshal to parse a stream of bytes with XML into a data structure. If the XML is malformed or cannot be mapped onto Plant, a descriptive error will be returned.

    ```go
        var p Plant
        if err := xml.Unmarshal(out, &p); err != nil {
            panic(err)
        }
        fmt.Println(p)
    ```
* The parent>child>plant field tag tells the encoder to nest all plants under <parent><child>...

    ```go
        tomato := &Plant{Id: 81, Name: "Tomato"}
        tomato.Origin = []string{"Mexico", "California"}

        type Nesting struct {
            XMLName xml.Name `xml:"nesting"`
            Plants  []*Plant `xml:"parent>child>plant"`
        }

        nesting := &Nesting{}
        nesting.Plants = []*Plant{coffee, tomato}
        out, _ = xml.MarshalIndent(nesting, " ", "  ")
        fmt.Println(string(out))
    ```

### Time Formatting / Parsing

* Go supports time formatting and parsing via pattern-based layouts.
* Here’s a basic example of formatting a time according to `RFC3339`, using the corresponding layout constant.

    ```go
        p := fmt.Println
        t := time.Now()
        p(t.Format(time.RFC3339))
    ```
* Time parsing uses the same layout values as Format.

    ```go
        t1, e := time.Parse(
        time.RFC3339,
        "2012-11-01T22:08:41+00:00")
        p(t1)
    ```

* Format and Parse use example-based layouts. Usually you’ll use a constant from time for these layouts,
* But you can also supply `custom layouts`. Layouts must use the `reference time Mon Jan 2 15:04:05 MST 2006` to show the pattern with which to format/parse a given time/string.
* The example time must be exactly as shown: the year 2006, 15 for the hour, Monday for the day of the week, etc.

    ```go
        p(t.Format("3:04PM"))
        p(t.Format("Mon Jan _2 15:04:05 2006"))
        p(t.Format("2006-01-02T15:04:05.999999-07:00"))
        form := "3 04 PM"
        t2, e := time.Parse(form, "8 41 PM")
        p(t2)
    ```
* For purely numeric representations you can also use standard string formatting with the extracted components of the time value.

    ```go
        fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
            t.Year(), t.Month(), t.Day(),
            t.Hour(), t.Minute(), t.Second()
        )
    ```
* Parse will return an error on malformed input explaining the parsing problem.

    ```go
        ansic := "Mon Jan _2 15:04:05 2006"
        _, e = time.Parse(ansic, "8:41PM")
        p(e)
    ```

### Arrays

* In Go, an array is a numbered sequence of elements of a specific length.

* Here we create an array, that will hold exactly 5 ints. 
* The type of elements and length are both part of the array’s type. 
* By default an array is `zero-valued`, which for ints means 0s.

    ```go
        var a [5]int
        fmt.Println("emp:", a)
    ```

* We can set a value at an index using the array[index] = value syntax, 
* and get a value with array[index].

    ```go
        a[4] = 100
        fmt.Println("set:", a)
        fmt.Println("get:", a[4])
    ```

* The builtin `len` returns the length of an array.

    ```go
        fmt.Println("len:", len(a))
    ```

* Use this syntax to declare and initialize an array in one line.

    ```go
        b := [5]int{1, 2, 3, 4, 5}
        fmt.Println("dcl:", b)
    ```
* Array types are `one-dimensional`, but you can compose types to `build multi-dimensional` data structures.

    ```go
        var twoD [2][3]int
        for i := 0; i < 2; i++ {
            for j := 0; j < 3; j++ {
                twoD[i][j] = i + j
            }
        }
        fmt.Println("2d: ", twoD)
    ```

* Note that arrays appear in the form [v1 v2 v3 ...] when printed with fmt.Println.

* You’ll see slices much more often than arrays in typical Go. We’ll look at slices next.

### slices

* Slices are a key data type in Go, giving a more powerful interface to sequences than arrays.

* Unlike arrays,`slices are typed only by the elements they contain (not the number of elements)`
* `To create an empty slice with non-zero length, use the builtin make`.
*  Here we make a slice of strings of length 3 (initially zero-valued).

    ```go
        s := make([]string, 3)
        fmt.Println("emp:", s)
    ```
* We can set and get just like with arrays.

    ```go
        s[0] = "a"
        s[1] = "b"
        s[2] = "c"
        fmt.Println("set:", s)
        fmt.Println("get:", s[2])
    ```

* `len` returns the length of the slice as expected.

    ```go
        fmt.Println("len:", len(s))
    ```

* In addition to these basic operations, slices support several more that make them richer than arrays. One is the builtin `append`, `which returns a slice containing one or more new values`.
* Note that we need to accept a return value from append as we may get a new slice value.

    ```go
        s = append(s, "d")
        s = append(s, "e", "f")
        fmt.Println("apd:", s)
    ```

* `Slices can also be copy’d`. Here we create an empty slice c of the same length as s and copy into c from s.

    ```go
        c := make([]string, len(s))
        copy(c, s)
        fmt.Println("cpy:", c)
    ```

* Slices support a `“slice” operator` with the syntax `slice[low:high]`. 
* For example, this gets a slice of the elements s[2], s[3], and s[4].

    ```go
        l := s[2:5]
        fmt.Println("sl1:", l)
    ```
* This slices up to (but excluding) s[5].

    ```go
        l = s[:5]
        fmt.Println("sl2:", l)
    ```

* And this slices up from (and including) s[2].

    ```go
        l = s[2:]
        fmt.Println("sl3:", l)
    ```

* We can declare and initialize a variable for slice in a single line as well.

    ```go
        t := []string{"g", "h", "i"}
        fmt.Println("dcl:", t)
    ```

* Slices can be composed into multi-dimensional data structures. 
* The `length of the inner slices can vary`, unlike with multi-dimensional arrays.

    ```go
        twoD := make([][]int, 3)
        for i := 0; i < 3; i++ {
            innerLen := i + 1
            twoD[i] = make([]int, innerLen)
            for j := 0; j < innerLen; j++ {
                twoD[i][j] = i + j
            }
        }
        fmt.Println("2d: ", twoD)
    ```

* Note that while slices are different types than arrays, they are rendered similarly by fmt.Println.

* Blog on Slices : http://blog.golang.org/2011/01/go-slices-usage-and-internals.html

### Maps

* Maps are Go’s built-in `associative data type (sometimes called hashes or dicts in other languages)`.

* To create an empty map, use the builtin make: make(map[key-type]val-type).

    ```go
        m := make(map[string]int)
    ```
* Set key/value pairs using typical name[key] = val syntax.

    ```go
        m["k1"] = 7
        m["k2"] = 13
    ```

* Printing a map with e.g. fmt.Println will show all of its key/value pairs.

    ```go
        fmt.Println("map:", m)
    ```

* Get a value for a key with name[key].

    ```go
        v1 := m["k1"]
        fmt.Println("v1: ", v1)
    ```

* The builtin len returns the number of key/value pairs when called on a map.

    ```go
        fmt.Println("len:", len(m))
    ```

* The builtin delete removes key/value pairs from a map.

    ```go
        delete(m, "k2")
        fmt.Println("map:", m)
    ```

* The optional second return value when getting a value from a map indicates if the key was present in the map.
* This can be used to disambiguate between missing keys and keys with zero values like 0 or "".
* Here we didn’t need the value itself, so we ignored it with the blank identifier _.

    ```go
        _, prs := m["k2"]
        fmt.Println("prs:", prs)
    ```

* You can also declare and initialize a new map in the same line with this syntax.

    ```go
        n := map[string]int{"foo": 1, "bar": 2}
        fmt.Println("map:", n)
    ```

* Note that maps appear in the form map[k:v k:v] when printed with fmt.Println.

### Range

* range iterates over elements in a variety of data structures.
* Let’s see how to use range with some of the data structures

* Here we use range to sum the numbers in a slice. Arrays work like this too.

    ```go
        nums := []int{2, 3, 4}
        sum := 0
        for _, num := range nums {
            sum += num
        }
        fmt.Println("sum:", sum)
    ```

* `range on arrays and slices` provides both the `index and value` for each entry. 
* Above we didn’t need the index, so we ignored it with the blank identifier _.
*  Sometimes we actually want the indexes though.

    ```go
        for i, num := range nums {
            if num == 3 {
                fmt.Println("index:", i)
            }
        }
    ```

* range on map iterates over key/value pairs.

    ```go
        kvs := map[string]string{"a": "apple", "b": "banana"}
        for k, v := range kvs {
            fmt.Printf("%s -> %s\n", k, v)
        }
    ```

* range can also iterate over just the keys of a map.

    ```go
        for k := range kvs {
           fmt.Println("key:", k)
        }
    ```

* range on strings iterates over Unicode code points. The first value is the starting byte index of the rune and the second the rune itself.

    ```go
        for i, c := range "go" {
           fmt.Println(i, c)
        }
    ```

### Functions

* Functions are central in Go. We’ll learn about functions with a few different examples.
* Here’s a function that takes two ints and returns their sum as an int.

    ```go
        func plus(a int, b int) int {
            return a + b
        }
    ```

* Go requires explicit returns, i.e. it won’t automatically return the value of the last expression.

* When you have multiple consecutive parameters of the same type, you may omit the type name for the like-typed parameters up to the final parameter that declares the type.

    ```go
        func plusPlus(a, b, c int) int {
            return a + b + c
        }
    ```
* Call a function just as you’d expect, with name(args).

    ```go
        res := plus(1, 2)
        fmt.Println("1+2 =", res)
        res = plusPlus(1, 2, 3)
        fmt.Println("1+2+3 =", res)
    ```

### Time

* Go offers extensive support for times and durations; here are some examples.
* We’ll start by getting the current time.

    ```go
        p := fmt.Println


        now := time.Now()
        p(now) // 2021-06-10 08:05:44.457982841 +0800 +08 m=+0.000074231
    ```
* You can build a time struct by providing the year, month, day, etc. `Times are always associated with a Location`, i.e. time zone.

    ```go
        p := fmt.Println
        then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

        p(then) // 2009-11-17 20:34:58.651387237 +0000 UTC
    ```

* You can extract the various components of the time value as expected.

    ```go
        p(then.Year()) // 2019
        p(then.Month()) // November
        p(then.Day()) // 17
        p(then.Hour()) // 20
        p(then.Minute()) // 34
        p(then.Second()) // 58
        p(then.Nanosecond()) // 651387237
        p(then.Location()) // UTC
    ```

* The Monday-Sunday Weekday is also available.

    ```go
        p(then.Weekday()) // Tuesday
    ```

* These methods compare two times, testing if the first occurs before, after, or at the same time as the second, respectively.

    ```go
        p(then.Before(now)) // true
        p(then.After(now)) // false
        p(then.Equal(now)) // false
    ```

* The Sub methods returns a Duration representing the interval between two times.

    ```go
        diff := now.Sub(then) // 101331h30m45.806595604s
        p(diff)
    ```

* We can compute the length of the duration in various units.

    ```go
        p(diff.Hours()) // 101331.51272405434
        p(diff.Minutes()) // 6.07989076344326e+06
        p(diff.Seconds()) // 3.647934458065956e+08
        p(diff.Nanoseconds()) // 364793445806595604
    ```

* You can use Add to advance a time by a given duration, or with a - to move backwards by a duration.

    ```go
        p(then.Add(diff)) // 2021-06-10 00:05:44.457982841 +0000 UTC
        p(then.Add(-diff)) // 1998-04-27 17:04:12.844791633 +0000 UTC
    ```

### Epoch
* A common requirement in programs is getting the number of seconds, milliseconds, or nanoseconds since the Unix epoch. Here’s how to do it in Go.

* Use time.Now with Unix or UnixNano to get elapsed time since the Unix epoch in seconds or nanoseconds, respectively.
* Note that there is no UnixMillis, so to get the milliseconds since epoch you’ll need to manually divide from nanoseconds.

    ```go
        now := time.Now()
        secs := now.Unix()
        nanos := now.UnixNano()

        millis := nanos / 1000000
        fmt.Println(now) // 2021-06-10 08:13:27.194057698 +0800 +08 m=+0.000082269
        fmt.Println(secs) // 1623284007
        fmt.Println(millis) // 1623284007194
        fmt.Println(nanos) // 1623284007194057698
    ```

* You can also convert integer seconds or nanoseconds since the epoch into the corresponding time.

    ```go
        fmt.Println(time.Unix(secs, 0)) //2021-06-10 08:13:27 +0800 +08
        fmt.Println(time.Unix(0, nanos)) // 2021-06-10 08:13:27.194057698 +0800 +08
    ```
### Regular Expressions

* Go offers built-in support for regular expressions. Here are some examples of common regexp-related tasks in Go.

* This tests whether a pattern matches a string.
* we used a string pattern directly,
    ```go
        match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
        fmt.Println(match) // true
    ```

* but for other regexp tasks you’ll need to Compile an optimized Regexp struct.

    ```go
        r, _ := regexp.Compile("p([a-z]+)ch")
    ```

* Many methods are available on these structs. Here’s a match test like we saw earlier.

    ```go
        fmt.Println(r.MatchString("peach")) // true
    ```

* This finds the match for the regexp.

    ```go
        fmt.Println(r.FindString("peach punch")) // peach
    ```

* This also finds the first match but returns the start and end indexes for the match instead of the matching text.

    ```go
        fmt.Println(r.FindStringIndex("peach punch")) // [0 5]
    ```

* The Submatch variants include information about both the whole-pattern matches and the submatches within those matches. For example this will return information for both p([a-z]+)ch and ([a-z]+).

    ```go
        fmt.Println(r.FindStringSubmatch("peach punch")) // [peach ea]
    ```
* Similarly this will return information about the indexes of matches and submatches.

    ```go
        fmt.Println(r.FindStringSubmatchIndex("peach punch")) // [0 5 1 3]
    ```

* The All variants of these functions apply to all matches in the input, not just the first. For example to find all matches for a regexp.

    ```go
        fmt.Println(r.FindAllString("peach punch pinch", -1)) // [peach punch pinch]
    ```

* These All variants are available for the other functions we saw above as well.

    ```go
        fmt.Println(r.FindAllStringSubmatchIndex("peach punch pinch", -1)) [[0 5 1 3] [6 11 7 9] [12 17 13 15]]
    ```

* Providing a non-negative integer as the second argument to these functions will limit the number of matches.

    ```go
        fmt.Println(r.FindAllString("peach punch pinch", 2)) // [peach punch]
    ```

* Our examples above had string arguments and used names like MatchString. We can also provide []byte arguments and drop String from the function name.

    ```go
        fmt.Println(r.Match([]byte("peach"))) // true
    ```

* When creating global variables with regular expressions you can use the MustCompile variation of Compile. MustCompile panics instead of returning an error, which makes it safer to use for global variables.

    ```go
        r = regexp.MustCompile("p([a-z]+)ch")
        fmt.Println(r) // p([a-z]+)ch
    ```

* The regexp package can also be used to replace subsets of strings with other values.

    ```go
        fmt.Println(r.ReplaceAllString("a peach", "<fruit>")) // a <fruit>
    ```

* The Func variant allows you to transform matched text with a given function.

    ```go
        in := []byte("a peach")
        out := r.ReplaceAllFunc(in, bytes.ToUpper)
        fmt.Println(string(out)) // a PEACH
    ```

### String Functions
* The standard library’s strings package provides many useful string-related functions. Here are some examples to give you a sense of the package.

* We alias fmt.Println to a shorter name as we’ll use it a lot below.
* Here’s a sample of the functions available in strings. Since these are functions from the package, not methods on the string object itself, we need pass the string in question as the first argument to the function.
* You can find more functions in the [strings](http://golang.org/pkg/strings/) package docs.

    ```go
        var p = fmt.Println
        p("Contains:  ", s.Contains("test", "es")) //Contains:   true
        p("Count:     ", s.Count("test", "t")) //Count:      2
        p("HasPrefix: ", s.HasPrefix("test", "te")) //HasPrefix:  true
        p("HasSuffix: ", s.HasSuffix("test", "st")) //HasSuffix:  true
        p("Index:     ", s.Index("test", "e")) //Index:      1
        p("Join:      ", s.Join([]string{"a", "b"}, "-")) //Join:       a-b
        p("Repeat:    ", s.Repeat("a", 5)) //Repeat:     aaaaa
        p("Replace:   ", s.Replace("foo", "o", "0", -1)) //Replace:    f00
        p("Replace:   ", s.Replace("foo", "o", "0", 1)) //Replace:    f0o
        p("Split:     ", s.Split("a-b-c-d-e", "-")) //Split:      [a b c d e]
        p("ToLower:   ", s.ToLower("TEST")) //ToLower:    test
        p("ToUpper:   ", s.ToUpper("test")) //ToUpper:    TEST
        p()
    ```

* Not part of strings, but worth mentioning here, are the mechanisms for getting the length of a string in bytes and getting a byte by index.

    ```go
        p("Len: ", len("hello")) //Len:  5
        p("Char:", "hello"[1]) // Char: 101
    ```

* Note that len and indexing above work at the byte level. Go uses UTF-8 encoded strings, so this is often useful as-is. If you’re working with potentially multi-byte characters you’ll want to use encoding-aware operations. See strings, bytes, runes and characters in Go for more information.

### String Formatting

* Go offers excellent support for string formatting in the printf tradition. Here are some examples of common string formatting tasks.

    ```go
        type point struct {
           x, y int
        }
    ```

* Go offers several printing “verbs” designed to format general Go values. 
* For example, this prints an instance of our point struct.
    ```go
        p := point{1, 2}
        fmt.Printf("%v\n", p) // {1 2}
    ```
* If the value is a struct, the `%+v` variant will include the struct’s field names.
    ```go
        fmt.Printf("%+v\n", p) //{x:1 y:2}
    ```
* The `%#v` variant prints a Go syntax representation of the value, i.e. the source code snippet that would produce that value.
    ```go
        fmt.Printf("%#v\n", p) //main.point{x:1, y:2}
    ```
* To print the type of a value, use `%T`.
    ```go
        fmt.Printf("%T\n", p) //main.point
    ```
* Formatting booleans is straight-forward.
    ```go
        fmt.Printf("%t\n", true) //true
    ```
* There are many options for formatting integers. Use `%d` for standard, `base-10` formatting.
    ```go
        fmt.Printf("%d\n", 123) //123
    ```
* This prints a binary representation.
    ```go
        fmt.Printf("%b\n", 14) //1110
    ```
* This prints the character corresponding to the given integer.
    ```go
        fmt.Printf("%c\n", 33) //!
    ```
* `%x` provides hex encoding.
    ```go
        fmt.Printf("%x\n", 456) //1c8
    ```
* There are also several formatting options for floats. For basic decimal formatting use `%f`.
    ```go
        fmt.Printf("%f\n", 78.9) //78.900000
    ```
* `%e` and `%E` format the float in (slightly different versions of) scientific notation.

    ```go
        fmt.Printf("%e\n", 123400000.0) //1.234000e+08
        fmt.Printf("%E\n", 123400000.0) //1.234000E+08
    ```
* For basic string printing use `%s`.
    ```go
        fmt.Printf("%s\n", "\"string\"") //"string"
    ```
* To `double-quote strings` as in Go source, use `%q`.
    ```go
        fmt.Printf("%q\n", "\"string\"") //"\"string\""
    ```
* As with integers seen earlier, `%x` renders the string in `base-16`, with two output characters per byte of input.
    ```go
        fmt.Printf("%x\n", "hex this") //6865782074686973
    ```
* To print a representation of a pointer, use `%p`.
    ```go
        fmt.Printf("%p\n", &p) //0xc00009e000
    ```
* When formatting numbers you will often want to control the width and precision of the resulting figure. To specify the width of an integer, use a number after the % in the verb. `By default the result will be right-justified and padded with spaces`.
    ```go
        fmt.Printf("|%6d|%6d|\n", 12, 345) //|    12|   345|
    ```
* You can also specify the width of printed floats, though usually you’ll also want to restrict the decimal precision at the same time with the width.precision syntax.
    ```go
        fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45) //|  1.20|  3.45|
    ```
* To `left-justify`, use the `-` flag.
    ```go
        fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45) //|1.20  |3.45  |
    ```
* You may also want to control width when formatting strings, especially to ensure that they align in table-like output. For basic right-justified width.
    ```go
        fmt.Printf("|%6s|%6s|\n", "foo", "b") //|   foo|     b|
    ```
* To `left-justify` use the `-` flag as with numbers.
    ```go
        fmt.Printf("|%-6s|%-6s|\n", "foo", "b") //|foo   |b     |
    ```
* So far we’ve seen `Printf`, which prints the formatted string to `os.Stdout`. `Sprintf formats` and `returns a string without printing it anywhere`.
    ```go
        s := fmt.Sprintf("a %s", "string")
        fmt.Println(s) //a string
    ```
* You can format+print to io.Writers other than os.Stdout using Fprintf.
    ```go
        fmt.Fprintf(os.Stderr, "an %s\n", "error") //an error
    ```

### Multiple Return Values

* Go has built-in support for `multiple return values`.
* This feature is used often in idiomatic Go, for example to return both result and error values from a function.

* The `(int, int)` in this function signature shows that the function returns 2 ints.

    ```go
        func vals() (int, int) {
           return 3, 7
        }
    ```

* Here we use the 2 different return values from the call with multiple assignment.

    ```go
        a, b := vals()
        fmt.Println(a) // 3
        fmt.Println(b) // 7
    ```
* If you only want a subset of the returned values, use the blank identifier _.

    ```go
        _, c := vals()
        fmt.Println(c) // 7
    ```
### Variadic Functions

* Variadic functions can be called with `any number of trailing arguments`.
* For example, `fmt.Println is a common variadic function`.

* Here’s a function that will take an arbitrary number of ints as arguments.

    ```go
        func sum(nums ...int) {
            fmt.Print(nums, " ")
            total := 0
            for _, num := range nums {
                total += num
            }
            fmt.Println(total)
        }
    ```

* Variadic functions can be called in the usual way with individual arguments.

    ```go
        sum(1, 2) // [1 2] 3
        sum(1, 2, 3) // [1 2 3] 6
    ```

* If you already have `multiple args in a slice`, apply them to a `variadic function using func(slice...)` like this.

    ```go
        nums := []int{1, 2, 3, 4}
        sum(nums...) // [1 2 3 4] 10
    ```

### Pointers

* Go supports pointers, allowing you to pass references to values and records within your program.

* We’ll show how pointers work in contrast to values with 2 functions: zeroval and zeroptr. *
* zeroval has an int parameter, so `arguments will be passed to it by value`. 
* zeroval will get `a copy of ival distinct from the one in the calling function`.

    ```go
        func zeroval(ival int) {
            ival = 0
        }
    ```

* zeroptr in contrast has an `*int parameter, meaning that it takes an int pointer`. 
* The `*iptr` code in the function body then dereferences the pointer from its memory address to the current value at that address.
* Assigning a value to a dereferenced pointer changes the value at the referenced address.

    ```go
        func zeroptr(iptr *int) {
            *iptr = 0
        }
    ```
    ```go
        i := 1
        fmt.Println("initial:", i) // inintial : 1
        zeroval(i)
        fmt.Println("zeroval:", i) // zeroval: 1
    ```
* The `&i syntax gives the memory address of i`, i.e. `a pointer to i`.

    ```go
        zeroptr(&i)
        fmt.Println("zeroptr:", i) // zeroptr: 0
    ```

* Pointers can be printed too.

    ```go
        fmt.Println("pointer:", &i) // pointer: 0xc0000ba010
    ```
* zeroval doesn’t change the i in main, but zeroptr does because it has a reference to the memory address for that variable.

### Structs

* Go’s `structs are typed collections of fields`.
* They’re useful for `grouping data together to form records`.

* This person struct type has name and age fields.

    ```go
        type person struct {
            name string
            age  int
        }
    ```
* newPerson constructs a new person struct with the given name.

    ```go
        func newPerson(name string) *person {
            p := person{name: name}
            p.age = 42
            return &p
        }
    ```

* You can safely return a pointer to local variable as a local variable will survive the scope of the function.

* This syntax creates a new struct.

    ```go
        fmt.Println(person{"Bob", 20}) //{Bob 20}
    ```
* You can name the fields when initializing a struct.

    ```go
        fmt.Println(person{name: "Alice", age: 30}) // {Alice 30}
    ```
* Omitted fields will be zero-valued.

    ```go
        fmt.Println(person{name: "Fred"}) // {Fred 0}
    ```
* An & prefix yields a pointer to the struct.

    ```go
        fmt.Println(&person{name: "Ann", age: 40}) // &{Ann 40}
    ```
* It’s idiomatic to encapsulate new struct creation in constructor functions
    ```go
        fmt.Println(newPerson("Jon")) // &{Jon 42}
    ```
* Access struct fields with a dot.
    ```go
        s := person{name: "Sean", age: 50}
        fmt.Println(s.name) // Sean
    ```
* `You can also use dots with struct pointers - the pointers are automatically dereferenced`.
    ```go
        sp := &s
        fmt.Println(sp.age) //50
    ```
* Structs are mutable.
    ```go
        sp.age = 51
        fmt.Println(sp.age) // 51
    ```
### Methods

* Go supports methods defined on struct types.

    ```go
        type rect struct {
            width, height int
        }
    ```
* This area method has a receiver type of `*rect`.

    ```go
        func (r *rect) area() int {
            return r.width * r.height
        }
    ```
* `Methods can be defined for either pointer or value receiver types`.
*  Here’s an example of a value receiver.

    ```go
        func (r rect) perim() int {
            return 2*r.width + 2*r.height
        }
    ```
* Here we call the 2 methods defined for our struct.
    ```go
        r := rect{width: 10, height: 5}
        fmt.Println("area: ", r.area()) // area:  50
        fmt.Println("perim:", r.perim()) // perim: 30
    ```

* Go `automatically handles conversion between values and pointers for method calls`.
*  You may want to `use a pointer receiver type to avoid copying on method calls` or `to allow the method to mutate the receiving struct.`

    ```go
        rp := &r
        fmt.Println("area: ", rp.area()) // area:  50
        fmt.Println("perim:", rp.perim()) // perim: 30
    ```
### Recursion

* Go supports recursive functions. Here’s a classic factorial example.

* This `fact function` calls itself until it reaches the base case of fact(0).

    ```go
        func fact(n int) int {
            if n == 0 {
                return 1
            }
            return n * fact(n-1)
        }

        fmt.Println(fact(7)) // 5040
    ```

### Interfaces

* Interfaces are `named collections of method signatures`.
* Here’s a basic interface for geometric shapes.

    ```go
        type geometry interface {
            area() float64
            perim() float64
        }
    ```
* For our example we’ll implement this interface on rect and circle types.
    ```go
        type rect struct {
            width, height float64
        }
        type circle struct {
            radius float64
        }
    ```
* `To implement an interface in Go, we just need to implement all the methods in the interface`. 
* Here we implement geometry on rects.
    ```go
        func (r rect) area() float64 {
            return r.width * r.height
        }
        func (r rect) perim() float64 {
            return 2*r.width + 2*r.height
        }
    ```
* The implementation for circles.
    ```go
        func (c circle) area() float64 {
            return math.Pi * c.radius * c.radius
        }
        func (c circle) perim() float64 {
            return 2 * math.Pi * c.radius
        }
    ```
* `If a variable has an interface type, then we can call methods that are in the named interface`.
* Here’s a generic measure function taking advantage of this to work on any geometry.
    ```go
        func measure(g geometry) {
            fmt.Println(g)
            fmt.Println(g.area())
            fmt.Println(g.perim())
        }
    ```
* The circle and rect struct types both implement the geometry interface so we can use instances of these structs as arguments to measure.

    ```go
        r := rect{width: 3, height: 4}
        c := circle{radius: 5}
        measure(r) // {3 4}, 12, 14
        measure(c) // {5}, 78.53981633974483, 31.41592653589793
    ```

### Goroutines

* A `goroutine is a lightweight thread of execution`.

    ```go
        func f(from string) {
            for i := 0; i < 3; i++ {
                fmt.Println(from, ":", i)
            }
        }
    ```

* Suppose we have a function call f(s). Here’s how we’d call that in the usual way, `running it synchronously`.
    ```go
        f("direct")
    ```
* `To invoke this function in a goroutine, use go f(s). This new goroutine will execute concurrently with the calling one`.
    ```go
        go f("goroutine")
    ```
* `You can also start a goroutine for an anonymous function call`.
    ```go
        go func(msg string) {
            fmt.Println(msg)
        }("going")
    ```
* Our two function calls are running asynchronously in separate goroutines now. Wait for them to finish (`for a more robust approach, use a WaitGroup`).

    ```go
        time.Sleep(time.Second)
        fmt.Println("done")
    ```
* When we run this program, we see the output of the blocking call first, then the interleaved output of the two goroutines.
* This interleaving reflects the goroutines being run concurrently by the Go runtime.
    ```shell
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run goroutines.go 
        direct : 0
        direct : 1
        direct : 2
        going
        goroutine : 0
        goroutine : 1
        goroutine : 2
        done
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run goroutines.go 
        direct : 0
        direct : 1
        direct : 2
        goroutine : 0
        goroutine : 1
        goroutine : 2
        going
        done
    ```

### Channels
* `Channels are the pipes that connect concurrent goroutines`. 
* You can send values into channels from one goroutine and receive those values into another goroutine.

* Create a `new channel with make(chan val-type)`. Channels are typed by the values they convey.
    ```go
        messages := make(chan string)
    ```
* `Send a value into a channel using the channel <- syntax`.
*  Here we send "ping" to the messages channel we made above, from a new goroutine.
    ```go
        go func() { messages <- "ping" }()
    ```
* `The <-channel syntax receives a value from the channel`.
*  Here we’ll receive the "ping" message we sent above and print it out.
    ```go
        msg := <-messages
        fmt.Println(msg) // ping
    ```
* When we run the program the "ping" message is successfully passed from one goroutine to another via our channel.

* `By default sends and receives block until both the sender and receiver are ready`.
*  This property allowed us to wait at the end of our program for the "ping" message without having to use any other synchronization.

### Channel Buffering

* By default channels are unbuffered, meaning that they will only accept sends (chan <-) if there is a corresponding receive (<- chan) ready to receive the sent value. Buffered channels accept a limited number of values without a corresponding receiver for those values.

* Here we make a channel of strings buffering up to 2 values.
    ```go
        messages := make(chan string, 2)
    ```
* Because this channel is buffered, we can send these values into the channel without a corresponding concurrent receive.
    ```go
        messages <- "buffered"
        messages <- "channel"
    ```
* Later we can receive these two values as usual.
    ```go
        fmt.Println(<-messages) // buffered
        fmt.Println(<-messages) // channel
    ```
### Channel Synchronization
* We can `use channels to synchronize execution across goroutines`.
* Here’s an example of using a blocking receive to wait for a goroutine to finish.
* When waiting for `multiple goroutines to finish, you may prefer to use a WaitGroup`.

* This is the function we’ll run in a goroutine.
* The done channel will be used to notify another goroutine that this function’s work is done.
    ```go
        func worker(done chan bool) {
            fmt.Print("working...")
            time.Sleep(time.Second)
            fmt.Println("done")

            // Send a value to notify that we’re done.
            done <- true
        }
    ```
* Start a worker goroutine, giving it the channel to notify on.
    ```go
        done := make(chan bool, 1)
        go worker(done)
        // Block until we receive a notification from the worker on the channel.
        <-done
    ```
* If you removed the <- done line from this program, the program would exit before the worker even started.

### Channel Directions
* When using channels as function parameters, you can specify if a channel is meant to only send or receive values.
* This specificity increases the type-safety of the program.

* This ping function only accepts a channel for sending values.
* It would be a compile-time error to try to receive on this channel.
    ```go
        func ping(pings chan<- string, msg string) {
            pings <- msg
        }
    ```
* The pong function accepts one channel for receives (pings) and a second for sends (pongs).
    ```go
        func pong(pings <-chan string, pongs chan<- string) {
            msg := <-pings
            pongs <- msg
        }
    ```
    ```go
        pings := make(chan string, 1)
        pongs := make(chan string, 1)
        ping(pings, "passed message")
        pong(pings, pongs)
        fmt.Println(<-pongs)
    ```

### Select

* Go’s `select` lets you wait on multiple channel operations.
* Combining goroutines and channels with select is a powerful feature of Go.

* For our example we’ll select across two channels.

    ```go
        c1 := make(chan string)
        c2 := make(chan string)
    ```
* Each channel will receive a value after some amount of time, to simulate e.g. blocking RPC operations executing in concurrent goroutines.

    ```go
        go func() {
            time.Sleep(1 * time.Second)
            c1 <- "one"
        }()
        go func() {
            time.Sleep(2 * time.Second)
            c2 <- "two"
        }()
    ```
* We’ll use select to await both of these values simultaneously, printing each one as it arrives.

    ```go
        for i := 0; i < 2; i++ {
            select {
            case msg1 := <-c1:
                fmt.Println("received", msg1) // received one
            case msg2 := <-c2:
                fmt.Println("received", msg2) // received two
            }
        }
    ```
* We receive the values "one" and then "two" as expected.

* Note that the total execution time is only `~2 seconds` since both the 1 and 2 second Sleeps execute concurrently.

### Timeouts

* Timeouts are important for programs that connect to external resources or that otherwise need to bound execution time. 
* `Implementing timeouts in Go is easy and elegant thanks to channels and select`.

* For our example, suppose we’re executing an external call that returns its result on a channel `c1` after `2s`.
* Note that `the channel is buffered, so the send in the goroutine is nonblocking. This is a common pattern to prevent goroutine leaks in case the channel is never read`.
    ```go
        c1 := make(chan string, 1)
        go func() {
            time.Sleep(2 * time.Second)
            c1 <- "result 1"
        }()
    ```
* Here’s the select implementing a timeout.
* `res := <-c1` awaits the result and
* `<-time.After` awaits a value to be sent after the timeout of 1s.
* Since select proceeds with the first receive that’s ready, we’ll take the timeout case if the operation takes more than the allowed 1s.
    ```go
        select {
        case res := <-c1:
            fmt.Println(res)
        case <-time.After(1 * time.Second):
            fmt.Println("timeout 1")
        }
        // prints : timeout 1
    ```
* If we allow a longer timeout of 3s, then the receive from c2 will succeed and we’ll print the result.
    ```go
        c2 := make(chan string, 1)
        go func() {
            time.Sleep(2 * time.Second)
            c2 <- "result 2"
        }()
        select {
        case res := <-c2:
            fmt.Println(res)
        case <-time.After(3 * time.Second):
            fmt.Println("timeout 2")
        }
        // prints : result 2
    ```
* Running this program shows the `first operation timing out` and the `second succeeding`.

### Non-Blocking Channel Operations

* `Basic sends and receives on channels are blocking`. 
* However, we `can use select with a default clause to implement non-blocking sends, receives, and even non-blocking multi-way selects`.
    ```go
        messages := make(chan string)
        signals := make(chan bool)
    ```
* Here’s a non-blocking receive. If a value is available on messages then select will take the <-messages case with that value. If not it will immediately take the default case.
    ```go
        select {
        case msg := <-messages:
            fmt.Println("received message", msg)
        default:
            fmt.Println("no message received")
        }
        // prints : no message received
    ```
* A non-blocking send works similarly. Here msg cannot be sent to the messages channel, because the channel has no buffer and there is no receiver. Therefore the default case is selected.
    ```go
        msg := "hi"
        select {
        case messages <- msg:
            fmt.Println("sent message", msg)
        default:
            fmt.Println("no message sent")
        }
        // prints : no message sent
    ```
* We can use multiple cases above the default clause to implement a multi-way non-blocking select. Here we attempt non-blocking receives on both messages and signals.
    ```go
        select {
        case msg := <-messages:
            fmt.Println("received message", msg)
        case sig := <-signals:
            fmt.Println("received signal", sig)
        default:
            fmt.Println("no activity")
        }
        // prints : no activity
    ```

### Closing Channels
* Closing a channel indicates that no more values will be sent on it.
* This can be useful to communicate completion to the channel’s receivers.

* In this example we’ll use a jobs channel to communicate work to be done from the main() goroutine to a worker goroutine.
* When we have no more jobs for the worker we’ll close the jobs channel.
    ```go
        jobs := make(chan int, 5)
        done := make(chan bool)
    ```
* Here’s the worker goroutine. It repeatedly receives from jobs with `j, more := <-jobs`.
* In this special 2-value form of receive, `the more value will be false if jobs has been closed` and all values in the channel have already been received. We use this to notify on done when we’ve worked all our jobs.
    ```go
        go func() {
            for {
                j, more := <-jobs
                if more {
                    fmt.Println("received job", j)
                } else {
                    fmt.Println("received all jobs")
                    done <- true
                    return
                }
            }
        }()
    ```
* This sends 3 jobs to the worker over the jobs channel, then closes it.
    ```go
        for j := 1; j <= 3; j++ {
            jobs <- j
            fmt.Println("sent job", j)
        }
        close(jobs)
        fmt.Println("sent all jobs")
    ```
* We await the worker using the synchronization approach we saw earlier.
    ```go
        <-done
    ```
### Range over Channels
* Earlier we saw how `for` and `range` provide iteration over basic data structures.
* We can also use this syntax to iterate over values received from a channel.

* We’ll iterate over 2 values in the queue channel.
    ```go
        queue := make(chan string, 2)
        queue <- "one"
        queue <- "two"
        close(queue)
    ```
* This range iterates over each element as it’s received from queue. Because we closed the channel above, the iteration terminates after receiving the 2 elements.
    ```go
        for elem := range queue {
        fmt.Println(elem)
    }
    // prints : one two
    ```
}
* This example also showed that it’s possible to close a non-empty channel but still have the remaining values be received.

### Timers
* We often want to execute Go code at some point in the future, or repeatedly at some interval.
* Go’s built-in `timer and ticker` features make both of these tasks easy.

* Timers represent a single event in the future. You tell the timer how long you want to wait, and it `provides a channel that will be notified at that time`.
* This timer will wait 2 seconds.
    ```go
        timer1 := time.NewTimer(2 * time.Second)
    ```
* The `<-timer1.C` blocks on the timer’s channel C until it sends a value indicating that the timer fired.
    ```go
        <-timer1.C
        fmt.Println("Timer 1 fired")
        // prints : Timer 1 fired
    ```
* If you just wanted to wait, `you could have used time.Sleep`.
* One reason a `timer may be useful is that you can cancel the timer before it fires`.
* Here’s an example of that.
    ```go
        timer2 := time.NewTimer(time.Second)
        go func() {
            <-timer2.C
            fmt.Println("Timer 2 fired")
        }()
        stop2 := timer2.Stop()
        if stop2 {
            fmt.Println("Timer 2 stopped")
        }
        // prints : Timer 2 stopped
    ```
* Give the timer2 enough time to fire, if it ever was going to, to show it is in fact stopped.
    ```go
        time.Sleep(2 * time.Second)
    ```
* The first timer will fire ~2s after we start the program, but the second should be stopped before it has a chance to fire.

### Tickers
* `Timers` are for when you want to do something `once in the future` -
* `Tickers` are for when you want to do `something repeatedly at regular intervals`.
* Here’s an example of a ticker that ticks periodically until we stop it.

* Tickers use a similar mechanism to timers: a channel that is sent values.
* Here we’ll use the select builtin on the channel to await the values as they arrive every 500ms.
    ```go
        ticker := time.NewTicker(500 * time.Millisecond)
        done := make(chan bool)
        go func() {
            for {
                select {
                case <-done:
                    return
                case t := <-ticker.C:
                    fmt.Println("Tick at", t)
                }
            }
        }()
        // prints : Tick at 2021-06-13 10:39:19.984730076 +0800 +08 m=+0.500913251
        // prints : Tick at 2021-06-13 10:39:20.484335584 +0800 +08 m=+1.000518768
        // prints : Tick at 2021-06-13 10:39:20.984986453 +0800 +08 m=+1.501169622
    ```
* Tickers can be stopped like timers. Once a ticker is stopped it won’t receive any more values on its channel. We’ll stop ours after 1600ms.
    ```go
        time.Sleep(1600 * time.Millisecond)
        ticker.Stop()
        done <- true
        fmt.Println("Ticker stopped")
        // prints : Ticker stopped
    ```
* When we run this program the ticker should tick 3 times before we stop it.

