package main

import (
	"time"
)

var conf Config

func cleanKeys() {
	now := time.Now()
	max := time.Duration(conf.MaxTime) * time.Second

	for key, _ := range conf.keys {
		if now.Sub(conf.keys[key]) > max {
			delete(conf.keys, key)
		}
	}
}

func scrape() {
	scrapePastes()
	scrapeGists()
}

func main() {
	conf = newConfig()
	initDB()

	for {
		scrape()
		time.Sleep(time.Duration(conf.Sleep) * time.Second)
		cleanKeys()
	}
}
