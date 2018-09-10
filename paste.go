package main

import (
	"encoding/json"
	"log"
	"time"
)

type Paste struct {
	ScrapeUrl string `json:"scrape_url"`
	Url       string `json:"full_url"`
	Key       string
	Content   string
}

func (p *Paste) Download() {
	_, exists := conf.keys[p.Key]
	if exists {
		return
	}

	log.Printf("[+] Downloading paste: %s\n", p.Key)

	resp := get(p.ScrapeUrl)
	p.Content = string(resp)
	conf.keys[p.Key] = time.Now()
}

func scrapePastes(c chan<- *ProcessItem) {
	var pastes []*Paste

	log.Println("[+] Checking for new pastes.")

	resp := get("https://scrape.pastebin.com/api_scraping.php?limit=100")
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

		item := &ProcessItem{Source: "Pastebin", Key: p.Key, Content: p.Content}
		c <- item
	}
}
