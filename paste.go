package main

import (
	"fmt"
	"time"
	"bytes"
)

type Paste struct {
	ScrapeUrl string `json:"scrape_url"`
	Url       string `json:"full_url"`
	Date      string
	Key       string
	Size      string
	Expire    string
	Title     string
	Syntax    string
	User      string
	Error     string
	Content   string
}

func (p *Paste) Download() {
	_, exists := conf.keys[p.Key]
	if exists {
		// fmt.Println("Already fetched this paste.")
		return
	}

	resp := get(p.ScrapeUrl)
    p.Content = string(resp)
	conf.keys[p.Key] = time.Now()
}

func (p *Paste) String() string {
	var b bytes.Buffer
	rule := "-----------------"

    b.WriteString(fmt.Sprintf("%s\n", rule))
	b.WriteString(fmt.Sprintf("Link: %s\n", p.Url))
	b.WriteString(fmt.Sprintf("Posted: %s\n", p.Date))
	b.WriteString(fmt.Sprintf("Expires: %s\n", p.Expire))
	b.WriteString(fmt.Sprintf("User: %s\n", p.User))
	b.WriteString(fmt.Sprintf("%s\n\n", rule))
	b.WriteString(fmt.Sprintf("%s", p.Content))

    return b.String()
}
