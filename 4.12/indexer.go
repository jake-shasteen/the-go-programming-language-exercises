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
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	"unicode"
)

type XKCDJSON struct {
	Month int `json:",string"`
	Num int 
	Link string
	Year int `json:",string"`
	News string
	SafeTitle string `json:"safe_title"`
	Transcript string
	Alt string
	Img string
	Title string
	Day int `json:",string"`
}

func main () {
	replacer := strings.NewReplacer(",", " ", "\n", " ", "\"", "", "'", "", "(", "", ")", "", "[", "", "]", "", "<", "", ">", "", "{", "", "}", "", ".", "", "?", "", "!", "", ";", "", ":", "", "-", "", "_", "", "/", "", "*", "")
	// searchIndex := make(map[string]([]int))
	comicIndex := make(map[int]*XKCDJSON)

	dir, err := os.Open("xkcd")
	handle(err)
	fileNames, err := dir.Readdirnames(0)
	handle(err)
	dir.Close()

	for _, fileName := range fileNames {
		fmt.Printf("%s read\n", fileName)

		file, err := os.Open("xkcd/" + fileName)
		handle(err)

		b, err := ioutil.ReadAll(file)
		handle(err)

		file.Close()

		index, err := strconv.Atoi(strings.Split(fileName, ".")[0])
		handle(err)

		comicIndex[index] = new(XKCDJSON) // allocate space for XKCDJSON
		err = json.Unmarshal(b, comicIndex[index])
		handle(err)

	}

	for index, json := range comicIndex {
		fmt.Printf("%d: %s %s \n", index, (*json).SafeTitle, removeSpace(replacer.Replace((*json).Transcript)))
	}
}

func handle(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func removeSpace(s string) string {
		var word bool
    rr := make([]rune, 0, len(s))
    for _, r := range s {
        if !unicode.IsSpace(r) {
            rr = append(rr, r)
            word = true
        } else if word {
        	word = false
        	rr = append(rr, ' ')
        }
    }
    return string(rr)
}