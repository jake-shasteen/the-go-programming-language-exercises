package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ChannelResponse struct {
	Number  int
	Payload []byte
}

func main() {
	ch := make(chan ChannelResponse)
	number := 1
	for number <= 2338 {
		go fetch(number, ch)
		number++
	}
	for number = 1; number <= 2338; number++ {
		if cr, ok := <-ch; ok {
			err := ioutil.WriteFile("xkcd2/"+strconv.Itoa(cr.Number)+".json", cr.Payload, 0644)
			if err != nil {
				fmt.Printf("Something went wrong when Writing\n")
			}
		} else {
			fmt.Printf("Something went wrong when receiving from channel\n")
		}
	}
	fmt.Printf("All done!\n")
	os.Exit(0)
}

func fetch(number int, ch chan<- ChannelResponse) {
	url := xkcdURL(number)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Something went wrong when Getting\n", number)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("Something went wrong when Reading\n", number)
	}
	ch <- ChannelResponse{Number: number, Payload: b}
}

func xkcdURL(number int) string {
	return "https://xkcd.com/" + strconv.Itoa(number) + "/info.0.json"
}
