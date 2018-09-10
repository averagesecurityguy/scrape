package main

import (
	"encoding/json"
	"log"
	"time"
)

type GistFile struct {
	Name    string `json:"filename"`
	Type    string `json:"language"`
	Url     string `json:"raw_url"`
	Size    int    `json:"size"`
	Content string `json:"content"`
}

type Gist struct {
	Url   string
	Date  string `json:"updated_at"`
	Key   string `json:"id"`
	User  string
	files []*GistFile
}

func (g *Gist) Download() {
	_, exists := conf.keys[g.Key]
	if exists {
		return
	}

	log.Printf("[+] Downloading gist: %s\n", g.Key)

	var gist map[string]*json.RawMessage
	data := getGithub(g.Url)

	err := json.Unmarshal(data, &gist)
	if err != nil {
		log.Printf("[-] Could not parse gist %s: %s\n.", g.Key, err.Error())
		return
	}

	// Decode each file object out of the Files map.
	var files map[string]*GistFile
	err = json.Unmarshal(*gist["files"], &files)
	if err != nil {
		log.Printf("[-] Could not parse gist file: %s\n", err.Error())
	}

	for k := range files {
		g.files = append(g.files, files[k])
	}

	conf.keys[g.Key] = time.Now()
}

func scrapeGists(c chan<- *ProcessItem) {
	var gists []*Gist

	log.Println("[+] Checking for new gists.")

	resp := getGithub("https://api.github.com/gists/public?per_page=100")
	err := json.Unmarshal(resp, &gists)
	if err != nil {
		log.Println("[-] Could not parse list of gists.")
		log.Printf("[-] %s.\n", err.Error())
		log.Println(string(resp))
		return
	}

	for i, _ := range gists {
		g := gists[i]
		g.Download()

		for _, gist := range g.files {
			item := &ProcessItem{Source: "Gist", Key: g.Key, Content: gist.Content}
			c <- item
		}
	}
}
