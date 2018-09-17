package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/boltdb/bolt"
)

type Keyword struct {
	Keyword string
	Bucket  string
}

type Regex struct {
	Regex    string
	compiled *regexp.Regexp
	Bucket   string
	Match    string
}

type Config struct {
	keys          map[string]time.Time
	db            *bolt.DB
	Keywords      []*Keyword // A list of keywords to search for in the data.
	Regexes       []*Regex   // A list of regular expressions to test against data.
	DbFile        string     `json:"database_file"` // File to use for the Store database.
	MaxSize       int        `json:"max_size"`      // Do not save files larger than this many bytes.
	MaxTime       int        `json:"max_time"`      // Max time, in seconds, to store previously downloaded keys.
	Sleep         int        // Time, in seconds, to wait between each run.
	GhToken       string     `json:"github_token"` // Github API token
	Save          bool
	LocalPath     string `json:"local_path"`
	FileBatchSize int    `json:"file_batch_size"`
	CertFile      string `json:"cert_file"`
	KeyFile       string `json:"key_file"`
	WebBatchSize     int `json:"web_batch_size"` // How many items to display per web page.
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

	if c.FileBatchSize == 0 {
		c.FileBatchSize = 100
	}

	if c.WebBatchSize == 0 {
		c.WebBatchSize = 25
	}

	return c
}
