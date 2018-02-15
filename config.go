package main

import (
	"log"
	"regexp"
	"time"
)

type Keyword struct {
	regex  *regexp.Regexp
	prefix string
}

type Config struct {
	keys     map[string]time.Time
	keywords []*Keyword
	ds       *KVStore
	maxSize  int           // Do not save files larger than this.
	maxTime  time.Duration // Max time to store previously downloaded keys.
	sleep    time.Duration // Time to wait between each run.
}

func newConfig() Config {
	var c Config

	c.keys = make(map[string]time.Time)
	c.maxSize = 100 * 1024 * 1024
	c.maxTime = 3600 * time.Second
	c.sleep = 60 * time.Second

	// Build our includes list.
	c.keywords = loadKeywords()

	// Create our connection to the key value store.
	ds, err := NewKVStore("data/scrape.db")
	if err != nil {
		log.Fatalf("[-] Cannot open database: %s\n", err)
	}

	c.ds = ds

	return c
}

func loadKeywords() []*Keyword {
	return []*Keyword{
		&Keyword{regexp.MustCompile("(?i)BEGIN PRIVATE KEY"), "privkey"},
		&Keyword{regexp.MustCompile("(?i)BEGIN DSA PRIVATE KEY"), "privkey"},
		&Keyword{regexp.MustCompile("(?i)BEGIN RSA PRIVATE KEY"), "privkey"},
		&Keyword{regexp.MustCompile("(?i)FULLZ"), "carder"},
		&Keyword{regexp.MustCompile("(?i)`password`"), "sqlpass"},
		&Keyword{regexp.MustCompile("(?i)proof of concept"), "exploit"},
		&Keyword{regexp.MustCompile("(?i)remote code execution"), "exploit"},
		&Keyword{regexp.MustCompile("AKIA[A-Z0-9]{16}"), "awskey"},
		&Keyword{regexp.MustCompile("\\$[0-9]\\$[a-zA-Z0-9]\\$[a-zA-Z0-9./=]+"), "pwhash"},
		&Keyword{regexp.MustCompile("[a-zA-Z0-9]+::[a-zA-Z0-9]{10}:[a-z0-9]{32}:[a-z0-9-]+"), "pwhash"},
		&Keyword{regexp.MustCompile("[a-zA-Z0-9-_]+:[0-9]+:[a-z0-9]{32}:[a-z0-9]{32}"), "pwhash"},
	}
}
