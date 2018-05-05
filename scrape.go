package main

import (
	"log"
	"time"

	"github.com/asggo/store"
)

var conf Config

func getStoreConn() *store.Store {
	var ds *store.Store
	var err error

	for tries := 1; tries < 20; tries += 2 {
		// Create our connection to the key value store.
		ds, err = store.NewStore(conf.DbFile)
		if err != nil {
			log.Printf("[-] Cannot open database: %s\n", err)
			time.Sleep(1 << uint(tries) * time.Millisecond)
		} else {
			break
		}
	}

	return ds
}

func initStore(s *store.Store) {
	for _, bucket := range conf.Buckets {
		s.CreateBucket(bucket)
	}
}

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
	scrapeDumpz()
}

func main() {
	conf = newConfig()

	// Ensure our database is initialized.
	ds := getStoreConn()
	initStore(ds)
	ds.Close()

	for {
		conf.ds = getStoreConn()
		scrape()
		conf.ds.Close()
		time.Sleep(time.Duration(conf.Sleep) * time.Second)
		cleanKeys()
	}
}
