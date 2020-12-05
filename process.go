package main

import (
	"encoding/json"
	"log"
	"strings"
)

type SaveItem struct {
	Location string
	Content string
}

func (s *SaveItem) Json() []byte {
	data, err := json.Marshal(s)
	if err != nil {
		log.Printf("[-] Could not create JSON for %s", s.Location)
		return nil
	}

	return data
}

type ProcessItem struct {
	Source  string
	Key     string
	Location string
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

	s := SaveItem{Location: p.Location, Content: p.Content}

	db.Write("pastes", p.Key, s.Json())
}

func (p *ProcessItem) Regexes() {
	for i := range conf.Regexes {
		r := conf.Regexes[i]

		switch r.Match {
		case "all":
			items := r.compiled.FindAllString(p.Content, -1)

			if items != nil {
				p.Save = true
				data := strings.Join(items, "\n")
				s := SaveItem{Location: p.Location, Content: data}

				db.Write(r.Bucket, p.Key, s.Json())
			}

			// for k := range items {
			// 	rKey := fmt.Sprintf("%s-%d", p.Source, k)
			// 	db.Write(r.Bucket, rKey, []byte(items[k]))
			// }
		case "one":
			match := r.compiled.FindString(p.Content)

			if match != "" {
				p.Save = true
				s := SaveItem{Location: p.Location, Content: match}

				db.Write(r.Bucket, p.Key, s.Json())
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
			db.Write(kwd.Bucket, p.Location, nil)
		}
	}
}

func process(sem chan struct{}, pi *ProcessItem) {
	sem <- struct{}{}

	pi.Regexes()
	pi.Keywords()
	pi.Write()

	<-sem
}
