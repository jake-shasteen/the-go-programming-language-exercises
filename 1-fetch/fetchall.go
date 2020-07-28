package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
}

func fetch(url string, ch chan <- string) {
	resp, err := http.Get(url)
	if err != nil {
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		os.Exit(1)
	}
	ch <- fmt.Sprintf("%s\n", b)
}