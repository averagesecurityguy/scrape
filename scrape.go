package main

import (
	"time"
)

var conf Config
var db *Database

func cleanKeys() {
	now := time.Now()
	max := time.Duration(conf.MaxTime) * time.Second

	for key, _ := range conf.keys {
		if now.Sub(conf.keys[key]) > max {
			delete(conf.keys, key)
		}
	}
}

func initDatabase() {
	db.CreateBucket("pastes")

	for _, kw := range conf.Keywords {
		db.CreateBucket(kw.Bucket)
	}

	for _, re := range conf.Regexes {
		db.CreateBucket(re.Bucket)
	}
}

func scrape(piChan chan<- *ProcessItem) {
	for {
		scrapeGithubEvents(piChan)
		scrapePastes(piChan)
		scrapeGists(piChan)
		scrapeFiles(piChan)

		time.Sleep(time.Duration(conf.Sleep) * time.Second)
		cleanKeys()
	}
}

func main() {
	conf = newConfig()
	db = newDatabase(conf.DbFile)

	if db != nil {
		initDatabase()

		go startWebServer()

		processItemChan := make(chan *ProcessItem, 100)
		processSemaphore := make(chan struct{}, 10)

		go scrape(processItemChan)

		for item := range processItemChan {
			go process(processSemaphore, item)
		}
	}
}
