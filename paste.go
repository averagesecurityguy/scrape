package main

import (
	"fmt"
	"log"
	"time"
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
	// Find and save specific data.
	if processCredentials(p.Content, p.Key) || processEmails(p.Content, p.Key) ||
		processPrivKey(p.Content, p.Key) || processAWSKeys(p.Content, p.Key) {
		conf.ds.Write("pastes", p.Key, []byte(p.Content))
	}

	// Save pastes that match any of our keywords. First match wins. Use these
	// to find interesting data that will eventually be processed with a more
	// specific method.
	save := false
	for i, _ := range conf.keywords {
		kwd := conf.keywords[i]
		key := fmt.Sprintf("%s-%s", kwd.prefix, p.Key)
		match := kwd.regex.FindString(p.Content)

		if match != "" {
			save = true
			conf.ds.Write("keywords", key, nil)
		}
	}

	if save {
		conf.ds.Write("pastes", p.Key, []byte(p.Content))
	}
}
