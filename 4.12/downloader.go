package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	number := 1
	for number <= 2338 {
		url := xkcdURL(number)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Something went wrong when Getting", number)
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Something went wrong when Reading", number)
		}
		filePtr, err := os.Create("xkcd/" + strconv.Itoa(number) + ".json")
		if err != nil {
			fmt.Printf("Something went wrong when Opening", number)
		}

		defer filePtr.Close()
		_, err = filePtr.Write(b)
		if err != nil {
			fmt.Printf("Something went wrong when writing", number)
		}
		filePtr.Sync()
		number++
	}
}

func xkcdURL(number int) string {
	return "https://xkcd.com/" + strconv.Itoa(number) + "/info.0.json"
}
