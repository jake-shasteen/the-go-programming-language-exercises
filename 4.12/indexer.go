// for each file in the folder
// load the contents and Unmarshal it into a proper struct.
// Add each struct to a map where the key is the comic number and the value is the resulting struct

// For each JSON object, add each word of the title, and each word of the transcript to a map
// where the key is a term and the value is a slice of comic numbers

// Make sure the word is lowercased and that punctuation : , [ ] ( ) ' " . < > etc is removed.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type XKCDJSON struct {
	Month      int `json:",string"`
	Num        int
	Link       string
	Year       int `json:",string"`
	News       string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        int `json:",string"`
}

func NewSearcher(dirPath string) func(search string) ([]*XKCDJSON, error) {
	replacer := strings.NewReplacer(",", " ", "\n", " ", "\"", "", "'", "", "(", "", ")", "", "[", "", "]", "", "<", "", ">", "", "{", "", "}", "", ".", "", "?", "", "!", "", ";", "", ":", "", "-", "", "_", "", "/", "", "*", "")
	searchIndex := make(map[string]([]int))
	comicIndex := make(map[int]*XKCDJSON)

	dir, err := os.Open(dirPath)
	handle(err)
	fileNames, err := dir.Readdirnames(0)
	handle(err)
	dir.Close()

	for _, fileName := range fileNames {
		file, err := os.Open(dirPath + "/" + fileName)
		handle(err)

		b, err := ioutil.ReadAll(file)
		handle(err)

		file.Close()

		index, err := strconv.Atoi(strings.Split(fileName, ".")[0])
		handle(err)
		if index != 404 {
			comicIndex[index] = new(XKCDJSON) // allocate space for XKCDJSON
			err = json.Unmarshal(b, comicIndex[index])
			if(err != nil) {
				fmt.Printf("error unmarshalling: %d\n", index)
			}
		}
	}

	for index, json := range comicIndex {

		uniqued := unique(strings.Split(removeSpace(replacer.Replace((*json).SafeTitle + " " + (*json).Transcript)), " "))

		for _, str := range uniqued {
			searchIndex[str] = append(searchIndex[str], index)
		}

	}

	return func (search string) ([]*XKCDJSON, error) {
		var results []*XKCDJSON
		if comicIndices, ok := searchIndex[search]; ok {
			for _, index := range comicIndices {
				results = append(results, comicIndex[index])
			}
		} else {
			return results, errors.New("Search term not found in index")
		}
		return results, nil
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

func unique(s []string) []string {
	set := make(map[string]struct{})
	for _, str := range s {
		set[strings.ToLower(str)] = struct{}{}
	}

	var result []string
	for str := range set {
		result = append(result, str)
	}
	return result
}
