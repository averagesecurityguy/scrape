package main

import (
	"log"
	"time"
	"encoding/json"
)

type Paste struct {
	ScrapeUrl string `json:"scrape_url"`
	Url       string `json:"full_url"`
	Date      string
	Key       string
	Size      int `json:",string"`
	Expire    int `json:",string"`
	Title     string
	Syntax    string
	User      string
	Error     string
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

func (p *Paste) Process() {
	processContent(p.Key, p.Content)
}

func scrapePastes() {
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
