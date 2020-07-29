// for each file in the folder
	// load the contents and Unmarshal it into a proper struct.
// Add each struct to a map where the key is the comic number and the value is the resulting struct

// For each JSON object, add each word of the title, and each word of the transcript to a map
	// where the key is a term and the value is a slice of comic numbers

// Make sure the word is lowercased and that punctuation : , [ ] ( ) ' " . < > etc is removed.


package main

import (
	"os"
	"fmt"
)

func main () {
	dir, err := os.Open("xkcd")
	names, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for _, name := range names {
		fmt.Printf("%s\n", name)
	}
}