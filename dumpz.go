package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type DumpList struct {
	Dumps []*DumpId
}

type DumpId struct {
	Id int
}

type Dump struct {
	Url     string
	Date    string
	Key     string
	Comment string
	Syntax  string `json:"lexer"`
	User    string
	Content string
}

func (d *Dump) Download() {
	_, exists := conf.keys[d.Key]
	if exists {
		return
	}

	url := fmt.Sprintf("https://dumpz.org/api/dump/%s", d.Key)

	log.Printf("[+] Downloading dumpz: %s\n", d.Key)

	data := get(url)
	err := json.Unmarshal(data, &d)
	if err != nil {
		log.Printf("[-] Could not parse dumpz %s: %s\n.", d.Key, err.Error())
		return
	}

	conf.keys[d.Key] = time.Now()
}

func (d *Dump) Process() {
	processContent(d.Key, d.Content)
}

func scrapeDumpz() {
	var ids DumpList

	log.Println("[+] Checking for new dumpz.")

	resp := get("https://dumpz.org/api/recent?limit=100&public=1")
	err := json.Unmarshal(resp, &ids)
	if err != nil {
		log.Println("[-] Could not parse list of dumpz.")
		log.Printf("[-] %s.\n", err.Error())
		log.Println(string(resp))
		return
	}

	for i, _ := range ids.Dumps {
		d := new(Dump)
		key := strconv.Itoa(ids.Dumps[i].Id)

		d.Key = key
		d.Download()
		d.Process()
	}
}
