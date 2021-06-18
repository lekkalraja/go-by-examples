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
    ```s
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

### Worker Pools
* we’ll look at how to implement a worker pool using goroutines and channels.
* Here’s the worker, of which we’ll run several concurrent instances.
* These workers will receive work on the jobs channel and send the corresponding results on results channel.
* We’ll sleep a second per job to simulate an expensive task.

    ```go
        func worker(id int, jobs <-chan int, results chan<- int) {
            for j := range jobs {
                fmt.Println("worker", id, "started  job", j)
                time.Sleep(time.Second)
                fmt.Println("worker", id, "finished job", j)
                results <- j * 2
            }
        }
    ```
* In order to use our pool of workers we need to send them work and collect their results.
* We make 2 channels for this.
    ```go
        jobs := make(chan int, numJobs)
        results := make(chan int, numJobs)
    ```
* This starts up 3 workers, initially blocked because there are no jobs yet.
    ```go
        for w := 1; w <= 3; w++ {
            go worker(w, jobs, results)
        }
    ```
* Here we send 5 jobs and then close that channel to indicate that’s all the work we have.
    ```go
        const numJobs = 5
        for j := 1; j <= numJobs; j++ {
            jobs <- j
        }
        close(jobs)
    ```
* Finally we collect all the results of the work.
* This also ensures that the worker goroutines have finished.
* `An alternative way to wait for multiple goroutines is to use a WaitGroup`.
    ```go
        for a := 1; a <= numJobs; a++ {
            <-results
        }
    ```
* Our running program shows the 5 jobs being executed by various workers. The program only takes about 2 seconds despite doing about 5 seconds of total work because there are 3 workers operating concurrently.

### WaitGroups
* To `wait for multiple goroutines to finish, we can use a wait group`.
* This is the function we’ll run in every goroutine.
* Note that `a WaitGroup must be passed to functions by pointer`.
* On return, notify the WaitGroup that we’re done.
    ```go
        func worker(id int, wg *sync.WaitGroup) {
            defer wg.Done()
            fmt.Printf("Worker %d starting\n", id)
            // Sleep to simulate an expensive task.
            time.Sleep(time.Second)
            fmt.Printf("Worker %d done\n", id)
        }
    ```
* This WaitGroup is used to wait for all the goroutines launched here to finish.
* Launch several goroutines and increment the WaitGroup counter for each.
    ```go
        var wg sync.WaitGroup
        for i := 1; i <= 5; i++ {
            wg.Add(1)
            go worker(i, &wg)
        }
    ```
* `Block until the WaitGroup counter goes back to 0`; all the workers notified they’re done.
    ```go
        wg.Wait()
    ```

* The order of workers starting up and finishing is likely to be different for each invocation.

### Rate Limiting
* `Rate limiting` is an important mechanism for `controlling resource utilization and maintaining quality of service`.
* `Go elegantly supports rate limiting with goroutines, channels, and tickers`.

* First we’ll look at basic rate limiting. Suppose we want to limit our handling of incoming requests.
* We’ll serve these requests off a channel of the same name.
    ```go
        requests := make(chan int, 5)
        for i := 1; i <= 5; i++ {
            requests <- i
        }
        close(requests)
    ```
* This limiter channel will receive a value every 200 milliseconds.
* This is the regulator in our rate limiting scheme.
    ```go
        limiter := time.Tick(200 * time.Millisecond)
    ```
* By blocking on a receive from the limiter channel before serving each request, we limit ourselves to 1 request every 200 milliseconds.
    ```go
        for req := range requests {
            <-limiter
            fmt.Println("request", req, time.Now())
        }
    ```
* We may want to allow short bursts of requests in our rate limiting scheme while preserving the overall rate limit.
* We can accomplish this by buffering our limiter channel. This burstyLimiter channel will allow bursts of up to 3 events.
    ```go
        burstyLimiter := make(chan time.Time, 3)
    ```
* Fill up the channel to represent allowed bursting.
    ```go
        for i := 0; i < 3; i++ {
            burstyLimiter <- time.Now()
        }
    ```
* Every 200 milliseconds we’ll try to add a new value to burstyLimiter, up to its limit of 3.
    ```go
        go func() {
            for t := range time.Tick(200 * time.Millisecond) {
                burstyLimiter <- t
            }
        }()
    ```
* Now simulate 5 more incoming requests. The first 3 of these will benefit from the burst capability of burstyLimiter.
    ```go
        burstyRequests := make(chan int, 5)
        for i := 1; i <= 5; i++ {
            burstyRequests <- i
        }
        close(burstyRequests)
        for req := range burstyRequests {
            <-burstyLimiter
            fmt.Println("request", req, time.Now())
        }
    ```
}
* Running our program we see the first batch of requests handled once every ~200 milliseconds as desired.
* For the second batch of requests we serve the first 3 immediately because of the burstable rate limiting, then serve the remaining 2 with ~200ms delays each.

### Atomic Counters
* The primary mechanism for managing state in Go is communication over channels.
* We saw this for example with worker pools. There are a few other options for managing state though.
* Here we’ll look at using the `sync/atomic` package for atomic counters accessed by multiple goroutines.

* We’ll use an unsigned integer to represent our (always-positive) counter.
    ```go
        var ops uint64
    ```
* A WaitGroup will help us wait for all goroutines to finish their work.
    ```go
        var wg sync.WaitGroup
    ```
* We’ll start 50 goroutines that each increment the counter exactly 1000 times.
* To atomically increment the counter we use `AddUint64`, giving it the `memory address of our ops counter with the & syntax`.
    ```go
        for i := 0; i < 50; i++ {
            wg.Add(1)
            go func() {
                for c := 0; c < 1000; c++ {
                    atomic.AddUint64(&ops, 1)
                }
                wg.Done()
            }()
        }
    ```
* Wait until all the goroutines are done.
    ```go
        wg.Wait()
    ```
* It’s safe to access ops now because we know no other goroutine is writing to it.
* Reading atomics safely while they are being updated is also possible, using functions like `atomic.LoadUint64`.
    ```go
        fmt.Println("ops:", ops)
        fmt.Println("ops:", atomic.LoadUnit64(&ops))
    ```
* We expect to get exactly 50,000 operations. Had we used the non-atomic ops++ to increment the counter, we’d likely get a different number, changing between runs, because the goroutines would interfere with each other. Moreover, `we’d get data race failures when running with the -race flag`.

### Mutexes

* In the previous example we saw how to manage simple counter state using atomic operations.
* For more complex state we can use a mutex to safely access data across multiple goroutines.

* For our example the state will be a map.
    ```go
        var state = make(map[int]int)
    ```
* This mutex will synchronize access to state.
    ```go
        var mutex = &sync.Mutex{}
    ```
* We’ll keep track of how many read and write operations we do.
    ```go
        var readOps uint64
        var writeOps uint64
    ```
* Here we start 100 goroutines to execute repeated reads against the state, once per millisecond in each goroutine.
* For each read we pick a key to access, 
  * `Lock()` the mutex to ensure exclusive access to the state,
  * `read` the value at the chosen key,
  * `Unlock()` the mutex, and increment the readOps count.
    ```go
        for r := 0; r < 100; r++ {
            go func() {
                total := 0
                for {
                    key := rand.Intn(5)
                    mutex.Lock()
                    total += state[key]
                    mutex.Unlock()
                    atomic.AddUint64(&readOps, 1)
                    time.Sleep(time.Millisecond)
                }
            }()
        }
    ```
* We’ll also start 10 goroutines to simulate writes, using the same pattern we did for reads.
    ```go
        for w := 0; w < 10; w++ {
            go func() {
                for {
                    key := rand.Intn(5)
                    val := rand.Intn(100)
                    mutex.Lock()
                    state[key] = val
                    mutex.Unlock()
                    atomic.AddUint64(&writeOps, 1)
                    time.Sleep(time.Millisecond)
                }
            }()
        }
    ```
* Let the 10 goroutines work on the state and mutex for a second.
    ```go
        time.Sleep(time.Second)
    ```
* Take and report final operation counts.
    ```go
        readOpsFinal := atomic.LoadUint64(&readOps)
        fmt.Println("readOps:", readOpsFinal)
        writeOpsFinal := atomic.LoadUint64(&writeOps)
        fmt.Println("writeOps:", writeOpsFinal)
    ```
* With a final lock of state, show how it ended up.
    ```go
        mutex.Lock()
        fmt.Println("state:", state)
        mutex.Unlock()
    ```
* Running the program shows that we executed about 90,000 total operations against our mutex-synchronized state.

    ```s
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run mutexs.go 
        readOps: 75400
        writeOps: 7527
        state: map[0:29 1:48 2:8 3:61 4:33]
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run mutexs.go 
        readOps: 75823
        writeOps: 7547
        state: map[0:81 1:68 2:39 3:42 4:78]
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run mutexs.go 
        readOps: 76288
        writeOps: 7639
        state: map[0:91 1:31 2:77 3:91 4:92]
    ```

### Stateful Goroutines

* In the previous example we used explicit locking with mutexes to synchronize access to shared state across multiple goroutines.
* Another option is to use the built-in synchronization features of goroutines and channels to achieve the same result.
* This channel-based approach aligns with Go’s ideas of sharing memory by communicating and having each piece of data owned by exactly 1 goroutine.

* In this example our state will be owned by a single goroutine.
* This will guarantee that the data is never corrupted with concurrent access.
* In order to read or write that state, other goroutines will send messages to the owning goroutine and receive corresponding replies.
* These readOp and writeOp structs encapsulate those requests and a way for the owning goroutine to respond.
    ```go
        type readOp struct {
            key  int
            resp chan int
        }
        type writeOp struct {
            key  int
            val  int
            resp chan bool
        }
    ```
* As before we’ll count how many operations we perform.
    ```go
        var readOps uint64
        var writeOps uint64
    ```
* The reads and writes channels will be used by other goroutines to issue read and write requests, respectively.
    ```go
        reads := make(chan readOp)
        writes := make(chan writeOp)
    ```
* Here is the goroutine that owns the state, which is a map as in the previous example but now private to the stateful goroutine.
* This goroutine repeatedly selects on the reads and writes channels, responding to requests as they arrive.
* A response is executed by first performing the requested operation and then sending a value on the response channel resp to indicate success (and the desired value in the case of reads).
    ```go
        go func() {
            var state = make(map[int]int)
            for {
                select {
                case read := <-reads:
                    read.resp <- state[read.key]
                case write := <-writes:
                    state[write.key] = write.val
                    write.resp <- true
                }
            }
        }()
    ```
* This starts 100 goroutines to issue reads to the state-owning goroutine via the reads channel.
* Each read requires constructing a readOp, sending it over the reads channel, and the receiving the result over the provided resp channel.
    ```go
        for r := 0; r < 100; r++ {
            go func() {
                for {
                    read := readOp{
                        key:  rand.Intn(5),
                        resp: make(chan int)}
                    reads <- read
                    <-read.resp
                    atomic.AddUint64(&readOps, 1)
                    time.Sleep(time.Millisecond)
                }
            }()
        }
    ```
* We start 10 writes as well, using a similar approach.
    ```go
        for w := 0; w < 10; w++ {
            go func() {
                for {
                    write := writeOp{
                        key:  rand.Intn(5),
                        val:  rand.Intn(100),
                        resp: make(chan bool)}
                    writes <- write
                    <-write.resp
                    atomic.AddUint64(&writeOps, 1)
                    time.Sleep(time.Millisecond)
                }
            }()
        }
    ```
* Let the goroutines work for a second.
* Finally, capture and report the op counts
    ```go
        time.Sleep(time.Second)
        readOpsFinal := atomic.LoadUint64(&readOps)
        fmt.Println("readOps:", readOpsFinal)
        writeOpsFinal := atomic.LoadUint64(&writeOps)
        fmt.Println("writeOps:", writeOpsFinal)
    ```
* For this particular case the goroutine-based approach was a bit more involved than the mutex-based one.
* It might be useful in certain cases though, for example where you have other channels involved or when managing multiple such mutexes would be error-prone. 
* You should use whichever approach feels most natural, especially with respect to understanding the correctness of your program.
    ```s
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run stateful_goroutines.go 
        readOps: 69319
        writeOps: 7025
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run stateful_goroutines.go 
        readOps: 74151
        writeOps: 7415
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run stateful_goroutines.go 
        readOps: 74034
        writeOps: 7438
    ```
### Command-Line Arguments
* Command-line arguments are a common way to parameterize execution of programs.
* For example, `go run hello.go uses run and hello.go arguments to the go program`.

* `os.Args` provides access to raw command-line arguments.
* Note that the `first value in this slice is the path to the program`,
* and os.Args[1:] holds the arguments to the program.
    ```go
        argsWithProg := os.Args
        argsWithoutProg := os.Args[1:]
    ```
* You can get individual args with normal indexing.
    ```go
        arg := os.Args[3]
        fmt.Println(argsWithProg) // [./command-line-arguments first second third]
        fmt.Println(argsWithoutProg) // [first second third]
        fmt.Println(arg) // third
    ```
* To experiment with command-line arguments it’s best to build a binary with go build first.
    ```s
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go build command-line-arguments.go 
        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ ./command-line-arguments first second third
        [./command-line-arguments first second third]
        [first second third]
        third
    ```
### Command-Line Flags
* Command-line flags are a common way to specify options for command-line programs.
* For example, `in wc -l the -l is a command-line flag`.

* Go provides a `flag package` supporting basic command-line flag parsing.
* We’ll use this package to implement our example command-line program.

* `Basic flag declarations are available for string, integer, and boolean options`.
* Here we declare a string flag word with a default value "foo" and a short description.
* This `flag.String function returns a string pointer (not a string value)`
    ```go
        wordPtr := flag.String("word", "foo", "a string")
    ```
* This declares numb and fork flags, using a similar approach to the word flag.
    ```go
        numbPtr := flag.Int("numb", 42, "an int")
        boolPtr := flag.Bool("fork", false, "a bool")
    ```
* It’s also possible to declare an option that uses an existing var declared elsewhere in the program.
* Note that we need to pass in a pointer to the flag declaration function.
    ```go
        var svar string
        flag.StringVar(&svar, "svar", "bar", "a string var")
    ```
* Once all flags are declared, call `flag.Parse()` to execute the command-line parsing.
    ```go
        flag.Parse()
    ```
* Here we’ll just dump out the parsed options and any trailing positional arguments.
* Note that we need to dereference the pointers with e.g. *wordPtr to get the actual option values.
    ```go
        fmt.Println("word:", *wordPtr)
        fmt.Println("numb:", *numbPtr)
        fmt.Println("fork:", *boolPtr)
        fmt.Println("svar:", svar)
        fmt.Println("tail:", flag.Args())
    ```
* To experiment with the command-line flags program it’s best to first compile it and then run the resulting binary directly.

* Try out the built program by first giving it values for all flags.
    ```s
        $ go build command-line-flags.go
        $ ./command-line-flags -word=opt -numb=7 -fork -svar=flag
        word: opt
        numb: 7
        fork: true
        svar: flag
        tail: []
    ```
* Note that if you omit flags they automatically take their default values.
    ```s
        $ ./command-line-flags -word=opt
        word: opt
        numb: 42
        fork: false
        svar: bar
        tail: []
    ```
* Trailing positional arguments can be provided after any flags.
    ```s
        $ ./command-line-flags -word=opt a1 a2 a3
        word: opt
        ...
        tail: [a1 a2 a3]
    ```
* Note that the flag package requires all flags to appear before positional arguments (otherwise the flags will be interpreted as positional arguments).
    ```s
        $ ./command-line-flags -word=opt a1 a2 a3 -numb=7
        word: opt
        numb: 42
        fork: false
        svar: bar
        tail: [a1 a2 a3 -numb=7]
    ```
* Use -h or --help flags to get automatically generated help text for the command-line program.
    ```s
        $ ./command-line-flags -h
        Usage of ./command-line-flags:
        -fork=false: a bool
        -numb=42: an int
        -svar="bar": a string var
        -word="foo": a string
    ```
* If you provide a flag that wasn’t specified to the flag package, the program will print an error message and show the help text again.
    ```s
        $ ./command-line-flags -wat
        flag provided but not defined: -wat
        Usage of ./command-line-flags:
    ```

### Command-Line Subcommands
* Some command-line tools, like the go tool or git have many subcommands, each with its own set of flags.
* For example, go build and go get are two different subcommands of the go tool.
* The flag package lets us easily define simple subcommands that have their own flags.

* We `declare a subcommand using the NewFlagSet function`, and proceed to define new flags specific for this subcommand.
    ```go
        fooCmd := flag.NewFlagSet("foo", flag.ExitOnError)
        fooEnable := fooCmd.Bool("enable", false, "enable")
        fooName := fooCmd.String("name", "", "name")
    ```
* For a different subcommand we can define different supported flags.
    ```go
        barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
        barLevel := barCmd.Int("level", 0, "level")
    ```
* The subcommand is expected as the first argument to the program.
    ```go
        if len(os.Args) < 2 {
            fmt.Println("expected 'foo' or 'bar' subcommands")
            os.Exit(1)
        }
    ```
* Check which subcommand is invoked.
* For every subcommand, we parse its own flags and have access to trailing positional arguments.
    ```go
        switch os.Args[1] {
        case "foo":
            fooCmd.Parse(os.Args[2:])
            fmt.Println("subcommand 'foo'")
            fmt.Println("  enable:", *fooEnable)
            fmt.Println("  name:", *fooName)
            fmt.Println("  tail:", fooCmd.Args())
        case "bar":
            barCmd.Parse(os.Args[2:])
            fmt.Println("subcommand 'bar'")
            fmt.Println("  level:", *barLevel)
            fmt.Println("  tail:", barCmd.Args())
        default:
            fmt.Println("expected 'foo' or 'bar' subcommands")
            os.Exit(1)
        }
    ```
* First invoke the foo subcommand.
    ```s
        $ go build command-line-subcommands.go
        $ ./command-line-subcommands foo -enable -name=joe a1 a2
        subcommand 'foo'
        enable: true
        name: joe
        tail: [a1 a2]
    ```
* Now try bar.
    ```s
        $ ./command-line-subcommands bar -level 8 a1
        subcommand 'bar'
        level: 8
        tail: [a1]
    ```
* But bar won’t accept foo’s flags.
    ```s
        $ ./command-line-subcommands bar -enable a1
        flag provided but not defined: -enable
        Usage of bar:
        -level int
                level
    ```

### Environment Variables
* Environment variables are a universal mechanism for conveying configuration information to Unix programs.
* Let’s look at how to set, get, and list environment variables.

* To `set a key/value pair, use os.Setenv`.
* To `get a value for a key, use os.Getenv`.
* This will return an empty string if the key isn’t present in the environment.
    ```go
        os.Setenv("FOO", "1")
        fmt.Println("FOO:", os.Getenv("FOO"))
        fmt.Println("BAR:", os.Getenv("BAR"))
    ```
* Use `os.Environ to list all key/value pairs in the environment`.
* This returns a slice of strings in the form KEY=value.
* You can strings.SplitN them to get the key and value. Here we print all the keys.
    ```go
        fmt.Println()
        for _, e := range os.Environ() {
            pair := strings.SplitN(e, "=", 2)
            fmt.Println(pair[0])
        }
    ```
}
    ```s
        $ go run environment-variables.go
        FOO: 1
        BAR: 

        SHELL
        SESSION_MANAGER
        QT_ACCESSIBILITY
        COLORTERM
        XDG_CONFIG_DIRS
        INTELLIJ_HOME
        XDG_MENU_PREFIX
        TERM_PROGRAM_VERSION
        GNOME_DESKTOP_SESSION_ID
        DERBY_HOME
        MANDATORY_PATH
        JAVA_HOME
        ....
    ```

### Sorting
* Go’s sort package implements sorting for builtins and user-defined types.
* We’ll look at sorting for builtins first.

* Sort methods are specific to the builtin type; here’s an example for strings.
* Note that `sorting is in-place, so it changes the given slice and doesn’t return a new one`.
    ```go
        strs := []string{"c", "a", "b"}
        sort.Strings(strs)
        fmt.Println("Strings:", strs)
    ```
* An example of sorting ints.
    ```go
        ints := []int{7, 2, 4}
        sort.Ints(ints)
        fmt.Println("Ints:   ", ints)
    ```
* We can also use sort to check if a slice is already in sorted order.
    ```go
        s := sort.IntsAreSorted(ints)
        fmt.Println("Sorted: ", s)
    ```

### Sorting by Functions

* Sometimes we’ll want to sort a collection by something other than its natural order.
* For example, suppose we wanted to sort strings by their length instead of alphabetically.
* Here’s an example of custom sorts in Go.

* In order to sort by a custom function in Go, we need a corresponding type.
* Here we’ve created a `byLength type` that is just an alias for the builtin []string type.
    ```go
        type byLength []string
    ```
* We implement `sort.Interface - Len, Less, and Swap` - on our type so we can use the sort package’s generic Sort function.
* Len and Swap will usually be similar across types and Less will hold the actual custom sorting logic.
* In our case we want to sort in order of increasing string length, so we use len(s[i]) and len(s[j]) here.
    ```go
        func (s byLength) Len() int {
            return len(s)
        }
        func (s byLength) Swap(i, j int) {
            s[i], s[j] = s[j], s[i]
        }
        func (s byLength) Less(i, j int) bool {
            return len(s[i]) < len(s[j])
        }
    ```
* With all of this in place, we can now implement our custom sort by converting the original fruits slice to byLength, and then use sort.Sort on that typed slice.
    ```go
        func main() {
            fruits := []string{"peach", "banana", "kiwi"}
            sort.Sort(byLength(fruits))
            fmt.Println(fruits) // [kiwi peach banana]
        }
    ```
* By following this same pattern of creating a custom type, implementing the three Interface methods on that type, and then calling sort.Sort on a collection of that custom type, we can sort Go slices by arbitrary functions.

### Collection Functions
* We often need our programs to perform operations on collections of data, like selecting all items that satisfy a given predicate or mapping all items to a new collection with a custom function.

* In some languages it’s idiomatic to use generic data structures and algorithms.
* Go does not support generics; in Go it’s common to provide collection functions if and when they are specifically needed for your program and data types.

* Here are some example collection functions for slices of strings.
* You can use these examples to build your own functions. Note that in some cases it may be clearest to just inline the collection-manipulating code directly, instead of creating and calling a helper function.

* Index returns the first index of the target string t, or -1 if no match is found.
    ```go
        func Index(vs []string, t string) int {
            for i, v := range vs {
                if v == t {
                    return i
                }
            }
            return -1
        }
    ```
* Include returns true if the target string t is in the slice.
    ```go
        func Include(vs []string, t string) bool {
            return Index(vs, t) >= 0
        }
    ```
* Any returns true if one of the strings in the slice satisfies the predicate f.
    ```go
        func Any(vs []string, f func(string) bool) bool {
            for _, v := range vs {
                if f(v) {
                    return true
                }
            }
            return false
        }
    ```
* All returns true if all of the strings in the slice satisfy the predicate f.
    ```go
        func All(vs []string, f func(string) bool) bool {
            for _, v := range vs {
                if !f(v) {
                    return false
                }
            }
            return true
        }
    ```
* Filter returns a new slice containing all strings in the slice that satisfy the predicate f.
    ```go
        func Filter(vs []string, f func(string) bool) []string {
            vsf := make([]string, 0)
            for _, v := range vs {
                if f(v) {
                    vsf = append(vsf, v)
                }
            }
            return vsf
        }
    ```
* Map returns a new slice containing the results of applying the function f to each string in the original slice.
    ```go
        func Map(vs []string, f func(string) string) []string {
            vsm := make([]string, len(vs))
            for i, v := range vs {
                vsm[i] = f(v)
            }
            return vsm
        }
    ```
* Here we try out our various collection functions.
    ```go
        var strs = []string{"peach", "apple", "pear", "plum"}
        fmt.Println(Index(strs, "pear")) // 2
        fmt.Println(Include(strs, "grape")) // false
        fmt.Println(Any(strs, func(v string) bool {
            return strings.HasPrefix(v, "p")
        })) // true
        fmt.Println(All(strs, func(v string) bool {
            return strings.HasPrefix(v, "p")
        })) // false
        fmt.Println(Filter(strs, func(v string) bool {
            return strings.Contains(v, "e")
        }))// [peach apple pear]
    ```
* The above examples all used anonymous functions, but you can also use named functions of the correct type.
    ```go
        fmt.Println(Map(strs, strings.ToUpper)) // [PEACH APPLE PEAR PLUM]
    ```
### Random Numbers
* Go’s math/rand package provides pseudorandom number generation.

* For example, `rand.Intn returns a random int n, 0 <= n < 100`.
    ```go
        fmt.Print(rand.Intn(100), ",")
        fmt.Print(rand.Intn(100))
        fmt.Println()
    ```
* `rand.Float64 returns a float64 f, 0.0 <= f < 1.0`.
    ```go
        fmt.Println(rand.Float64())
    ```
* This can be used to generate `random floats in other ranges, for example 5.0 <= f' < 10.0`.
    ```go
        fmt.Print((rand.Float64()*5)+5, ",")
        fmt.Print((rand.Float64() * 5) + 5)
        fmt.Println()
    ```
* The default number generator is deterministic, so it’ll produce the same sequence of numbers each time by default. To produce varying sequences, give it a seed that changes.
* `Note that this is not safe to use for random numbers you intend to be secret, use crypto/rand for those`.
    ```go
        s1 := rand.NewSource(time.Now().UnixNano())
        r1 := rand.New(s1)
    ```
* Call the resulting rand.Rand just like the functions on the rand package.
    ```go
        fmt.Print(r1.Intn(100), ",")
        fmt.Print(r1.Intn(100))
        fmt.Println()
    ```
* If you seed a source with the same number, it produces the same sequence of random numbers.
    ```go
        s2 := rand.NewSource(42)
        r2 := rand.New(s2)
        fmt.Print(r2.Intn(100), ",")
        fmt.Print(r2.Intn(100))
        fmt.Println()
        s3 := rand.NewSource(42)
        r3 := rand.New(s3)
        fmt.Print(r3.Intn(100), ",")
        fmt.Print(r3.Intn(100))
    ```
* See the `math/rand package docs` for references on other random quantities that Go can provide.

### Number Parsing
* Parsing numbers from strings is a basic but common task in many programs; here’s how to do it in Go.

* The `built-in package strconv provides the number parsing`.

* With `ParseFloat`, this 64 tells how many bits of precision to parse.
    ```go
        f, _ := strconv.ParseFloat("1.234", 64)
        fmt.Println(f)
    ```
* For ParseInt, the 0 means infer the base from the string. 64 requires that the result fit in 64 bits.
    ```go
        i, _ := strconv.ParseInt("123", 0, 64)
        fmt.Println(i)
    ```
* ParseInt will recognize hex-formatted numbers.
    ```go
        d, _ := strconv.ParseInt("0x1c8", 0, 64)
        fmt.Println(d)
    ```
* A ParseUint is also available.
    ```go
        u, _ := strconv.ParseUint("789", 0, 64)
        fmt.Println(u)
    ```
* Atoi is a convenience function for basic base-10 int parsing.
    ```go
        k, _ := strconv.Atoi("135")
        fmt.Println(k)
    ```
* Parse functions return an error on bad input.
    ```go
        _, e := strconv.Atoi("wat")
        fmt.Println(e)
    ```

### URL Parsing
* URLs provide a uniform way to locate resources. Here’s how to parse URLs in Go.

* We’ll parse this example URL, which includes a scheme, authentication info, host, port, path, query params, and query fragment.
    ```go
        s := "postgres://user:pass@host.com:5432/path?k=v#f"
    ```
* Parse the URL and ensure there are no errors.
    ```go
        u, err := url.Parse(s)
        if err != nil {
            panic(err)
        }
    ```
* Accessing the scheme is straightforward.
    ```go
        fmt.Println(u.Scheme)
    ```
* User contains all authentication info; call Username and Password on this for individual values.
    ```go
        fmt.Println(u.User)
        fmt.Println(u.User.Username())
        p, _ := u.User.Password()
        fmt.Println(p)
    ```
* The Host contains both the hostname and the port, if present. Use SplitHostPort to extract them.
    ```go
        fmt.Println(u.Host)
        host, port, _ := net.SplitHostPort(u.Host)
        fmt.Println(host)
        fmt.Println(port)
    ```
* Here we extract the path and the fragment after the #.
    ```go
        fmt.Println(u.Path)
        fmt.Println(u.Fragment)
    ```
* To get query params in a string of k=v format, use RawQuery.
* You can also parse query params into a map. The parsed query param maps are from strings to slices of strings, so index into [0] if you only want the first value.
    ```go
        fmt.Println(u.RawQuery)
        m, _ := url.ParseQuery(u.RawQuery)
        fmt.Println(m)
        fmt.Println(m["k"][0])
    ```

### SHA1 Hashes

* SHA1 hashes are frequently used to compute short identities for binary or text blobs.
* For example, the git revision control system uses SHA1s extensively to identify versioned files and directories. Here’s how to compute SHA1 hashes in Go.

* Go implements several hash functions in various `crypto/* packages`.
    ```go
        import (
            "crypto/sha1"
        )

        s := "sha1 this string"
    ```
* The pattern for generating a hash is `sha1.New()`, `sha1.Write(bytes)`, `then sha1.Sum([]byte{})`
* Here we start with a new hash.
    ```go
        h := sha1.New()
    ```
* Write expects bytes. If you have a string s, use []byte(s) to coerce it to bytes.
    ```go
        h.Write([]byte(s))
    ```
* This gets the finalized hash result as a byte slice. The argument to Sum can be used to append to an existing byte slice: it usually isn’t needed.
    ```fo
        bs := h.Sum(nil)
    ```
* SHA1 values are often printed in hex, for example in git commits.
* Use the %x format verb to convert a hash results to a hex string.
    ```go
        fmt.Println(s)
        fmt.Printf("%x\n", bs)
    ```
* You can compute other hashes using a similar pattern to the one shown above. For example, to compute MD5 hashes import crypto/md5 and use md5.New().

* Note that if you need cryptographically secure hashes, you should carefully research hash strength!

### Base64 Encoding
* Go provides built-in support for base64 encoding/decoding.

* This syntax imports the encoding/base64 package with the b64 name instead of the default base64. It’ll save us some space below.
    ```go
        import (
            b64 "encoding/base64"
            "fmt"
        )
    ```

* Here’s the string we’ll encode/decode.
    ```go
        data := "abc123!?$*&()'-=@~"
    ```
* Go supports both standard and URL-compatible base64.
* Here’s how to encode using the standard encoder.
* The encoder requires a []byte so we convert our string to that type.
    ```go
        sEnc := b64.StdEncoding.EncodeToString([]byte(data))
        fmt.Println(sEnc)
    ```
* Decoding may return an error, which you can check if you don’t already know the input to be well-formed.
    ```go
        sDec, _ := b64.StdEncoding.DecodeString(sEnc)
        fmt.Println(string(sDec))
    ```
* This encodes/decodes using a URL-compatible base64 format.
    ```go
        uEnc := b64.URLEncoding.EncodeToString([]byte(data))
        fmt.Println(uEnc)
        uDec, _ := b64.URLEncoding.DecodeString(uEnc)
        fmt.Println(string(uDec))
    ```
* The string encodes to slightly different values with the standard and URL base64 encoders (trailing + vs -) but they both decode to the original string as desired.
    ```s
        $ go run base64-encoding.go
        YWJjMTIzIT8kKiYoKSctPUB+
        abc123!?$*&()'-=@~
        YWJjMTIzIT8kKiYoKSctPUB-
        abc123!?$*&()'-=@~
    ```

### Reading Files
* Reading and writing files are basic tasks needed for many Go programs.
* First we’ll look at some examples of reading files.
* Reading files requires checking most calls for errors.
* This helper will streamline our error checks below.
    ```go
        func check(e error) {
            if e != nil {
                panic(e)
            }
        }
    ```

* Perhaps the most basic file reading task is slurping a file’s entire contents into memory.
    ```go
        dat, err := ioutil.ReadFile("errors.go")
        check(err)
        fmt.Print(string(dat))
    ```
* You’ll often want more control over how and what parts of a file are read. For these tasks, start by Opening a file to obtain an `os.File` value.
    ```go
        f, err := os.Open("/tmp/dat")
        check(err)
    ```
* Read some bytes from the beginning of the file. Allow up to 5 to be read but also note how many actually were read.
    ```go
        b1 := make([]byte, 5)
        n1, err := f.Read(b1)
        check(err)
        fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))
    ```
* You can also Seek to a known location in the file and Read from there.
    ```go
        o2, err := f.Seek(6, 0)
        check(err)
        b2 := make([]byte, 2)
        n2, err := f.Read(b2)
        check(err)
        fmt.Printf("%d bytes @ %d: ", n2, o2)
        fmt.Printf("%v\n", string(b2[:n2]))
    ```
* The io package provides some functions that may be helpful for file reading.
* For example, reads like the ones above can be more robustly implemented with ReadAtLeast.
    ```go
        o3, err := f.Seek(6, 0)
        check(err)
        b3 := make([]byte, 2)
        n3, err := io.ReadAtLeast(f, b3, 2)
        check(err)
        fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))
    ```
* There is no built-in rewind, but Seek(0, 0) accomplishes this.
    ```go
        _, err = f.Seek(0, 0)
        check(err)
    ```
* The bufio package implements a buffered reader that may be useful both for its efficiency with many small reads and because of the additional reading methods it provides.
    ```go
        r4 := bufio.NewReader(f)
        b4, err := r4.Peek(5)
        check(err)
        fmt.Printf("5 bytes: %s\n", string(b4))
    ```
* Close the file when you’re done (usually this would be scheduled immediately after Opening with defer).
    ```go
        f.Close()
    ```

### Writing Files
* Writing files in Go follows similar patterns to the ones we saw earlier for reading.
* To start, here’s how to dump a string (or just bytes) into a file.
    ```go
        d1 := []byte("hello\ngo\n")
        err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
        check(err)
    ```
* For more granular writes, open a file for writing.
    ```go
        f, err := os.Create("/tmp/dat2")
        check(err)
    ```
* It’s idiomatic to defer a Close immediately after opening a file.
    ```go
        defer f.Close()
    ```
* You can Write byte slices as you’d expect.
    ```go
        d2 := []byte{115, 111, 109, 101, 10}
        n2, err := f.Write(d2)
        check(err)
        fmt.Printf("wrote %d bytes\n", n2)
    ```
* A WriteString is also available.
    ```go
        n3, err := f.WriteString("writes\n")
        check(err)
        fmt.Printf("wrote %d bytes\n", n3)
    ```
* Issue a Sync to flush writes to stable storage.
    ```go
        f.Sync()
    ```
* `bufio` provides buffered writers in addition to the buffered readers we saw earlier.
    ```go
        w := bufio.NewWriter(f)
        n4, err := w.WriteString("buffered\n")
        check(err)
        fmt.Printf("wrote %d bytes\n", n4)
    ```
* Use Flush to ensure all buffered operations have been applied to the underlying writer.
    ```go
        w.Flush()
    ```

### Line Filters
* A line filter is a common type of program that reads input on stdin, processes it, and then prints some derived result to stdout. grep and sed are common line filters.

* Here’s an example line filter in Go that writes a capitalized version of all input text. You can use this pattern to write your own Go line filters.

* Wrapping the unbuffered os.Stdin with a buffered scanner gives us a convenient Scan method that advances the scanner to the next token; which is the next line in the default scanner.
    ```go
        scanner := bufio.NewScanner(os.Stdin)
    ```
* Text returns the current token, here the next line, from the input.
    ```go
        for scanner.Scan() {
            ucl := strings.ToUpper(scanner.Text())
            fmt.Println(ucl)
        }
    ```
* Check for errors during Scan. End of file is expected and not reported by Scan as an error.
    ```go
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "error:", err)
            os.Exit(1)
        }
    ```

### Directories
* Go has several useful functions for working with directories in the file system.
* Create a new sub-directory in the current working directory.
    ```go
        err := os.Mkdir("subdir", 0755)
        check(err)
    ```
* When creating temporary directories, it’s good practice to defer their removal. `os.RemoveAll` will delete a whole directory tree (`similarly to rm -rf`).
    ```go
        defer os.RemoveAll("subdir")
    ```
* Helper function to create a new empty file.
    ```go
        createEmptyFile := func(name string) {
            d := []byte("")
            check(ioutil.WriteFile(name, d, 0644))
        }
        createEmptyFile("subdir/file1")
    ```
* We can create a hierarchy of directories, `including parents with MkdirAll. This is similar to the command-line mkdir -p`.
    ```go
        err = os.MkdirAll("subdir/parent/child", 0755)
        check(err)
        createEmptyFile("subdir/parent/file2")
        createEmptyFile("subdir/parent/file3")
        createEmptyFile("subdir/parent/child/file4")
    ```
* ReadDir lists directory contents, returning a slice of os.FileInfo objects.
    ```go
        c, err := ioutil.ReadDir("subdir/parent")
        check(err)
        fmt.Println("Listing subdir/parent")
        for _, entry := range c {
            fmt.Println(" ", entry.Name(), entry.IsDir())
        }
    ```
* `Chdir` lets us change the current working directory, similarly to cd.
    ```go
        err = os.Chdir("subdir/parent/child")
        check(err)
    ```
* Now we’ll see the contents of subdir/parent/child when listing the current directory.
    ```go
        c, err = ioutil.ReadDir(".")
        check(err)
        fmt.Println("Listing subdir/parent/child")
        for _, entry := range c {
            fmt.Println(" ", entry.Name(), entry.IsDir())
        }
    ```
* cd back to where we started.
    ```go
        err = os.Chdir("../../..")
        check(err)
    ```
* We can also visit a directory recursively, including all its sub-directories.
* `Walk` accepts a callback function to handle every file or directory visited.
    ```go
        fmt.Println("Visiting subdir")
        err = filepath.Walk("subdir", visit)
    ```
* visit is called for every file or directory found recursively by `filepath.Walk`.
    ```go
        func visit(p string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            fmt.Println(" ", p, info.IsDir())
            return nil
        }
    ```

### Temporary Files and Directories
* Throughout program execution, we often want to create data that isn’t needed after the program exits. Temporary files and directories are useful for this purpose since they don’t pollute the file system over time.

* The easiest way to create a temporary file is by calling `ioutil.TempFile`. It creates a file and opens it for reading and writing. We provide "" as the first argument, so ioutil.TempFile will create the file in the default location for our OS.
    ```go
        f, err := ioutil.TempFile("", "sample")
        check(err)
    ```
* Display the name of the temporary file. On Unix-based OSes the directory will likely be /tmp. The file name starts with the prefix given as the second argument to ioutil.TempFile and the rest is chosen automatically to ensure that concurrent calls will always create different file names.
    ```go
        fmt.Println("Temp file name:", f.Name())
    ```
* Clean up the file after we’re done. The OS is likely to clean up temporary files by itself after some time, but it’s good practice to do this explicitly.
    ```go
        defer os.Remove(f.Name())
    ```
* We can write some data to the file.
    ```go
        _, err = f.Write([]byte{1, 2, 3, 4})
        check(err)
    ```
* If we intend to write many temporary files, we may prefer to create a temporary directory.
* ioutil.TempDir’s arguments are the same as TempFile’s, but it returns a directory name rather than an open file.
    ```go
        dname, err := ioutil.TempDir("", "sampledir")
        check(err)
        fmt.Println("Temp dir name:", dname)
        defer os.RemoveAll(dname)
    ```
* Now we can synthesize temporary file names by prefixing them with our temporary directory.
    ```go
        fname := filepath.Join(dname, "file1")
        err = ioutil.WriteFile(fname, []byte{1, 2}, 0666)
        check(err)
    ```

### Testing
* Unit testing is an important part of writing principled Go programs.
* The `testing package` provides the tools we need to write unit tests and the go test command runs tests.

* For the sake of demonstration, this code is in package main, but it could be any package.
* Testing code typically lives in the same package as the code it tests.

* We’ll be testing this simple implementation of an integer minimum. Typically, the code we’re testing would be in a source file named something like intutils.go, and the test file for it would then be named intutils_test.go.
    ```go
        func IntMin(a, b int) int {
            if a < b {
                return a
            }
            return b
        }
    ```
* A test is created by writing a function with a name beginning with Test.
* `t.Error*` will report test failures but continue executing the test.
* `t.Fatal*` will report test failures and stop the test immediately.

    ```go
        func TestIntMinBasic(t *testing.T) {
            ans := IntMin(2, -2)
            if ans != -2 {
                t.Errorf("IntMin(2, -2) = %d; want -2", ans)
            }
        }
    ```
* Writing tests can be repetitive, so it’s idiomatic to use a table-driven style, where test inputs and expected outputs are listed in a table and a single loop walks over them and performs the test logic.
* `t.Run` enables running “subtests”, one for each table entry. These are shown separately when executing `go test -v`.
    ```go
        func TestIntMinTableDriven(t *testing.T) {
            var tests = []struct {
                a, b int
                want int
            }{
                {0, 1, 0},
                {1, 0, 0},
                {2, -2, -2},
                {0, -1, -1},
                {-1, 0, -1},
            }

            for _, tt := range tests {
                testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
                t.Run(testname, func(t *testing.T) {
                    ans := IntMin(tt.a, tt.b)
                    if ans != tt.want {
                        t.Errorf("got %d, want %d", ans, tt.want)
                    }
                })
            }
        }
    ```
### If/Else
* Branching with if and else in Go is straight-forward.
    ```go
        if 7%2 == 0 {
            fmt.Println("7 is even")
        } else {
            fmt.Println("7 is odd")
        }
    ```
* You can have an if statement without an else.
    ```go
        if 8%4 == 0 {
            fmt.Println("8 is divisible by 4")
        }
    ```
* A statement can precede conditionals; any variables declared in this statement are available in all branches.
    ```go
        if num := 9; num < 0 {
            fmt.Println(num, "is negative")
        } else if num < 10 {
            fmt.Println(num, "has 1 digit")
        } else {
            fmt.Println(num, "has multiple digits")
        }
    ```
* Note that you don’t need parentheses around conditions in Go, but that the braces are required.
* There is `no ternary if in Go`, so you’ll need to use a full if statement even for basic conditions.

### For
* for is Go’s only looping construct. Here are some basic types of for loops.
* The most basic type, with a single condition.
    ```go
        i := 1
        for i <= 3 {
            fmt.Println(i)
            i = i + 1
        }
    ```
* A classic initial/condition/after for loop.
    ```go
        for j := 7; j <= 9; j++ {
            fmt.Println(j)
        }
    ```
* for without a condition will loop repeatedly until you break out of the loop or return from the enclosing function.
    ```go
        for {
            fmt.Println("loop")
            break
        }
    ```
* You can also continue to the next iteration of the loop.
    ```go
        for n := 0; n <= 5; n++ {
            if n%2 == 0 {
                continue
            }
            fmt.Println(n)
        }
    ```
### Variables
* In Go, variables are explicitly declared and used by the compiler to e.g. check type-correctness of function calls.

* var declares 1 or more variables.
    ```go
        var a = "initial"
        fmt.Println(a)
    ```
* You can declare multiple variables at once.
    ```go
        var b, c int = 1, 2
        fmt.Println(b, c)
    ```
* Go will infer the type of initialized variables.
    ```go
        var d = true
        fmt.Println(d)
    ```
* Variables declared without a corresponding initialization are zero-valued.
* For example, the zero value for an int is 0.
    ```go
        var e int
        fmt.Println(e)
    ```
* The `:=` syntax is shorthand for declaring and initializing a variable, e.g. for `var f string = "apple"` in this case.
    ```go
        f := "apple"
        fmt.Println(f)
    ```

### Values
* Go has various value types including strings, integers, floats, booleans, etc. Here are a few basic examples.

* Strings, which can be added together with +.
    ```go
        fmt.Println("go" + "lang")
    ```
* Integers and floats.
    ```go
        fmt.Println("1+1 =", 1+1)
        fmt.Println("7.0/3.0 =", 7.0/3.0)
    ```
* Booleans, with boolean operators as you’d expect.
    ```go
        fmt.Println(true && false)
        fmt.Println(true || false)
        fmt.Println(!true)
    ```

### GoodBye World
* Our Last program will print the classic "GoodBye world” message. Here’s the full source code.

    ```go
        package main
        import "fmt"
        func main() {
            fmt.Println("hello world")
        }
    ```