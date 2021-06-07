package main

import (
	"bufio"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://gobyexample.com/")

	if err != nil {
		log.Panicf("Something went wrong while fetching response : %v \n", err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Panicf("Something went wrong while Reading response : %v \n", err)
	}
}

/*
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$ go run http_client.go
	2021/06/07 09:45:19 <!DOCTYPE html>
	2021/06/07 09:45:19 <html>
	2021/06/07 09:45:19   <head>
	2021/06/07 09:45:19     <meta charset="utf-8">
	2021/06/07 09:45:19     <title>Go by Example</title>
	2021/06/07 09:45:19     <link rel=stylesheet href="site.css">
	2021/06/07 09:45:19   </head>
	2021/06/07 09:45:19   <body>
	2021/06/07 09:45:19     <div id="intro">
	2021/06/07 09:45:19       <h2><a href="./">Go by Example</a></h2>
	2021/06/07 09:45:19       <p>
	2021/06/07 09:45:19         <a href="http://golang.org">Go</a> is an
	2021/06/07 09:45:19         open source programming language designed for
	2021/06/07 09:45:19         building simple, fast, and reliable software.
	2021/06/07 09:45:19         Please read the official
	2021/06/07 09:45:19         <a href="https://golang.org/doc/tutorial/getting-started">documentation</a>
	2021/06/07 09:45:19         to learn a bit about Go code, tools packages,
	2021/06/07 09:45:19         and modules.
	2021/06/07 09:45:19       </p>
	2021/06/07 09:45:19
	2021/06/07 09:45:19       <p>
	2021/06/07 09:45:19         <em>Go by Example</em> is a hands-on introduction
	2021/06/07 09:45:19         to Go using annotated example programs. Check out
	2021/06/07 09:45:19         the <a href="hello-world">first example</a> or
	2021/06/07 09:45:19         browse the full list below.
	2021/06/07 09:45:19       </p>
	2021/06/07 09:45:19
	2021/06/07 09:45:19       <ul>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="hello-world">Hello World</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="values">Values</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="variables">Variables</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="constants">Constants</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="for">For</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="if-else">If/Else</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="switch">Switch</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="arrays">Arrays</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="slices">Slices</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="maps">Maps</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="range">Range</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="functions">Functions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="multiple-return-values">Multiple Return Values</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="variadic-functions">Variadic Functions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="closures">Closures</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="recursion">Recursion</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="pointers">Pointers</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="structs">Structs</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="methods">Methods</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="interfaces">Interfaces</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="errors">Errors</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="goroutines">Goroutines</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="channels">Channels</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="channel-buffering">Channel Buffering</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="channel-synchronization">Channel Synchronization</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="channel-directions">Channel Directions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="select">Select</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="timeouts">Timeouts</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="non-blocking-channel-operations">Non-Blocking Channel Operations</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="closing-channels">Closing Channels</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="range-over-channels">Range over Channels</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="timers">Timers</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="tickers">Tickers</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="worker-pools">Worker Pools</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="waitgroups">WaitGroups</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="rate-limiting">Rate Limiting</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="atomic-counters">Atomic Counters</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="mutexes">Mutexes</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="stateful-goroutines">Stateful Goroutines</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="sorting">Sorting</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="sorting-by-functions">Sorting by Functions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="panic">Panic</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="defer">Defer</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="collection-functions">Collection Functions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="string-functions">String Functions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="string-formatting">String Formatting</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="regular-expressions">Regular Expressions</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="json">JSON</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="xml">XML</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="time">Time</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="epoch">Epoch</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="time-formatting-parsing">Time Formatting / Parsing</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="random-numbers">Random Numbers</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="number-parsing">Number Parsing</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="url-parsing">URL Parsing</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="sha1-hashes">SHA1 Hashes</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="base64-encoding">Base64 Encoding</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="reading-files">Reading Files</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="writing-files">Writing Files</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="line-filters">Line Filters</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="file-paths">File Paths</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="directories">Directories</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="temporary-files-and-directories">Temporary Files and Directories</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="testing">Testing</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="command-line-arguments">Command-Line Arguments</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="command-line-flags">Command-Line Flags</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="command-line-subcommands">Command-Line Subcommands</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="environment-variables">Environment Variables</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="http-clients">HTTP Clients</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="http-servers">HTTP Servers</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="context">Context</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="spawning-processes">Spawning Processes</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="execing-processes">Exec'ing Processes</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="signals">Signals</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19         <li><a href="exit">Exit</a></li>
	2021/06/07 09:45:19
	2021/06/07 09:45:19       </ul>
	2021/06/07 09:45:19       <p class="footer">
	2021/06/07 09:45:19         by <a href="https://markmcgranaghan.com">Mark McGranaghan</a> | <a href="https://github.com/mmcgrana/gobyexample">source</a> | <a href="https://github.com/mmcgrana/gobyexample#license">license</a>
	2021/06/07 09:45:19       </p>
	2021/06/07 09:45:19     </div>
	2021/06/07 09:45:19   </body>
	2021/06/07 09:45:19 </html>
	raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-examples$
*/
