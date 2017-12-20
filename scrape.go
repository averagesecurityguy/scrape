package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var conf Config

func cleanKeys() {
	now := time.Now()

	for key, _ := range conf.keys {
		if now.Sub(conf.keys[key]) > conf.maxTime {
			delete(conf.keys, key)
		}
	}
}

func get(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not access %s\n.", url)
		return []byte("")
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Received HTTP error %d\n", resp.StatusCode)
		return []byte("")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Could not read response from %s\n.", url)
		return []byte("")
	}

	return body
}

func scrape() {
	var pastes []Paste

	resp := get("https://pastebin.com/api_scraping.php?limit=250&lang=test")
	err := json.Unmarshal(resp, &pastes)
	if err != nil {
		fmt.Println("Could not parse list of pastes.")
		return
	}

	for i, _ := range pastes {
		p := pastes[i]
		p.Download()
		process(p)
	}
}

func main() {
	conf = newConfig()
	for {
		scrape()
		time.Sleep(conf.sleep)
		cleanKeys()
	}
}
