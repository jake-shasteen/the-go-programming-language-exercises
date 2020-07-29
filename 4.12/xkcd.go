/*
Exercise 4.12:
	The popular web comic xkcd has a JSON interface. For example, a request to https://xkcd.com/571/info.0.json
	produces a detailed description of comic 571, one of many favorites.
	
	Download each URL (once!) and build an offline index.

	Write a tool xkcd that, using this index, prints the URL and transcript of each comic
	that matches a search term provided on the command line.
*/

// could have an in-memory index process running. Communciate through channels?
// This sounds like a separate process should run and host the index though.

package main

import (
	"os"
	"fmt"
)

func main() {

	argc := len(os.Args)

	if (argc < 2) {
		fmt.Printf("Usage: xkcd <search term>\n")
		os.Exit(1)
	}

	searchTerm := os.Args[1]

	search := NewSearcher("xkcd")

	results, err := search(searchTerm)
	if(err != nil) {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for _, comic := range results {
		fmt.Printf("%d: %s:\n %s\n", comic.Num, comic.SafeTitle, comic.Transcript)
	}
}
