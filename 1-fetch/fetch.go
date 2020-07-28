package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

/*
	os.Args
	os.Stdin
	os.Stdout
	os.Stderr
	os.Exit

	http.Get() (url: URL) => resp: Response
	Response {
		Body
		Status
		etc
	}

	fmt.Fprintf
	fmt.Printf

	ioutil.ReadAll -- makes a big buffer and puts all of the thing in instead of streaming.
	equiv to js:

	let result;

	eventStream
	.on("data", data => result.data += data)
	.on("end", console.log('you can use result now'))
*/

func main() {
	// args are urls
	for _, url := range os.Args[1:] { // slice off the 0th arg
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong when getting %s\n", url)
			os.Exit(1)
		}
		result, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong when buffering body of %s\n", url)
			os.Exit(1)
		}
		fmt.Printf("%s", result)
	}
}
