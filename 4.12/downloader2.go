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
	for number = 0; number <= 2338; number++ {
		cr, ok := <-ch
		if ok {
			err := ioutil.WriteFile("xkcd2/"+strconv.Itoa(cr.Number)+".json", cr.Payload, 0644)

			if err != nil {
				fmt.Printf("Something went wrong when Writing")
			}
		} else {
			fmt.Printf("Something went wrong when receiving from channel")
		}
	}
	os.Exit(0)
}

func fetch(number int, ch chan<- ChannelResponse) {
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
	ch <- ChannelResponse{Number: number, Payload: b}
}

func xkcdURL(number int) string {
	return "https://xkcd.com/" + strconv.Itoa(number) + "/info.0.json"
}
