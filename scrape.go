package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/asggo/store"
)

var conf Config

func getStoreConn() *store.Store {
	var ds *store.Store
	var err error

	for tries := 1; tries < 20; tries += 2 {
		// Create our connection to the key value store.
		ds, err = store.NewStore(conf.dbFile)
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
	s.CreateBucket("awskeys")
	s.CreateBucket("emails")
	s.CreateBucket("creds")
	s.CreateBucket("privkeys")
	s.CreateBucket("keywords")
	s.CreateBucket("pastes")
}

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
		log.Printf("[-] Could not access %s.\n", url)
		return []byte("")
	}

	if resp.StatusCode != 200 {
		log.Printf("[-] Received HTTP error %d.\n", resp.StatusCode)
		return []byte("")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[-] Could not read response from %s.\n", url)
		return []byte("")
	}

	return body
}

func scrape() {
	var pastes []*Paste

	log.Println("[+] Checking for new pastes.")

	resp := get("https://pastebin.com/api_scraping.php?limit=100")
	err := json.Unmarshal(resp, &pastes)
	if err != nil {
		log.Println("[-] Could not parse list of pastes.")
		log.Printf("[-] %s.\n", err.Error())
		log.Println(string(resp))
		return
	}

	for i, _ := range pastes {
		p := pastes[i]
		p.Download()
		p.Process()
	}
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
		time.Sleep(conf.sleep)
		cleanKeys()
	}
}
