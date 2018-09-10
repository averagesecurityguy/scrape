package main

import (
	"fmt"
	"log"
	"strings"
)

type ProcessItem struct {
	Source  string
	Key     string
	Content string
	Save    bool
}

func (p *ProcessItem) Write() {
	if (conf.Save == false) || (p.Save == false) {
		return
	}

	if len(p.Content) > conf.MaxSize {
		return
	}

	db.Write("pastes", p.Key, []byte(p.Content))
}

func (p *ProcessItem) Regexes() {
	for i, _ := range conf.Regexes {
		r := conf.Regexes[i]

		switch r.Match {
		case "all":
			items := r.compiled.FindAllString(p.Content, -1)

			if items != nil {
				p.Save = true
			}

			for k := range items {
				rKey := fmt.Sprintf("%s-%d", p.Key, k)
				db.Write(r.Bucket, rKey, []byte(items[k]))
			}
		case "one":
			match := r.compiled.FindString(p.Content)

			if match != "" {
				p.Save = true
				db.Write(r.Bucket, p.Key, []byte(match))
			}
		default:
		}
	}
}

func (p *ProcessItem) Keywords() {
	for i, _ := range conf.Keywords {
		kwd := conf.Keywords[i]

		if strings.Contains(strings.ToLower(p.Content), strings.ToLower(kwd.Keyword)) {
			p.Save = true
			db.Write(kwd.Bucket, p.Key, nil)
		}
	}
}

func process(sem chan struct{}, pi *ProcessItem) {
	log.Printf("[+] Processing %s:%s.\n", pi.Source, pi.Key)
	sem <- struct{}{}

	pi.Regexes()
	pi.Keywords()
	pi.Write()

	<-sem
}
