package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/asggo/store"
)

type Keyword struct {
	Keyword string
	Prefix  string
}

type Regex struct {
	Regex    string
	compiled *regexp.Regexp
	Prefix   string
	Match    string
}

type Config struct {
	keys     map[string]time.Time
	ds       *store.Store
	Keywords []*Keyword // A list of keywords to search for in the data.
	Regexes  []*Regex   // A list of regular expressions to test against data.
	Buckets  []string   `json:"buckets"`       // List of buckets we need to create.
	DbFile   string     `json:"database_file"` // File to use for the Store database.
	MaxSize  int        `json:"max_size"`      // Do not save files larger than this many bytes.
	MaxTime  int        `json:"max_time"`      // Max time, in seconds, to store previously downloaded keys.
	Sleep    int        // Time, in seconds, to wait between each run.
	GhToken  string     `json:"github_token"` // Github API token
	Save     bool
}

func newConfig() Config {
	var c Config

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("[-] Could not read config file.")
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal("[-] Could not parse config file.")
	}

	c.keys = make(map[string]time.Time)

	// Compile our regular expressions
	for i, _ := range c.Regexes {
		r := c.Regexes[i]
		r.compiled = regexp.MustCompile(r.Regex)
	}

	return c
}
