package main

import (
	"encoding/json"
	"log"
	"time"
)

type PushEventCommit struct {
	Key   string `json:"sha"`
	Url   string
	Files []GithubCommitFile
}

type PushEvent struct {
	Key     json.Number `json:"push_id"`
	Commits []PushEventCommit
}

type GithubEvent struct {
	Key     string `json:"id"`
	Type    string
	Date    string `json:"created_at"`
	Payload json.RawMessage
}

type GithubCommitFile struct {
	Url string `json:"raw_url"`
	Content string
}

type GithubCommit struct {
	Files []GithubCommitFile
}

func (p *PushEvent) Download() {
	_, exists := conf.keys[string(p.Key)]
	if exists {
		return
	}

	for i := range p.Commits {
		var ghc GithubCommit

		data := getGithub(p.Commits[i].Url)
		err := json.Unmarshal(data, &ghc)
		if err != nil {
			log.Printf("[-] Could not parse commit %s: %s\n.", p.Commits[i].Key, err.Error())
			continue
		}

		for j := range ghc.Files {
			f := ghc.Files[j]

			f.Content = string(get(f.Url))
			p.Commits[i].Files = append(p.Commits[i].Files, f)
		}
	}

	conf.keys[string(p.Key)] = time.Now()
}

func scrapeGithubEvents(c chan<- *ProcessItem) {
	var events []*GithubEvent

	log.Println("[+] Checking for new Github events.")

	resp := getGithub("https://api.github.com/events")
	err := json.Unmarshal(resp, &events)
	if err != nil {
		log.Println("[-] Could not parse list of events.")
		log.Printf("[-] %s.\n", err.Error())
		log.Println(string(resp))
		return
	}

	log.Printf("[+] Processing %d events.\n", len(events))

	for i := range events {
		if events[i].Type == "PushEvent" {
			var pe PushEvent

			err := json.Unmarshal(events[i].Payload, &pe)
			if err != nil {
				log.Printf("[-] Could not parse payload for %s\n", events[i].Key)
				log.Printf("[-] %s.\n", err.Error())
				log.Println(string(events[i].Payload))
				continue
			}

			pe.Download()
			for i := range pe.Commits {
				for j := range pe.Commits[i].Files {
					f := pe.Commits[i].Files[j]

					item := &ProcessItem{Source: "GithubCommit", Location: f.Url, Key: string(pe.Key), Content: f.Content}
					c <- item
				}
			}
		}
	}
}
