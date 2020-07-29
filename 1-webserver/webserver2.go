package main

import (
  "fmt"
  "log"
  "net/http"
  "sync"
)

var mu := sync.Mutex
var count int

func main() {
  http.HandleFunc("/", anyRoute)
  http.HandleFunc("/count", countRoute)
  log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func anyRoute(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
  fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func countRoute(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}