package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", anyRoute)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func anyRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}