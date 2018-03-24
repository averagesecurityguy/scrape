package main

import (
	"regexp"
	"time"

	"github.com/asggo/store"
)

type Keyword struct {
	regex  *regexp.Regexp
	prefix string
}

type Config struct {
	keys     map[string]time.Time
	keywords []*Keyword
	ds       *store.Store
	dbFile   string
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

	// Build our keyword list.
	c.keywords = loadKeywords()
	c.dbFile = "data/scrape.db"

	return c
}

func loadKeywords() []*Keyword {
	return []*Keyword{
		&Keyword{regexp.MustCompile("(?i)FULLZ"), "carder"},
		&Keyword{regexp.MustCompile("(?i)`password`"), "sqlpass"},
		&Keyword{regexp.MustCompile("(?i)proof of concept"), "exploit"},
		&Keyword{regexp.MustCompile("(?i)remote code execution"), "exploit"},
		&Keyword{regexp.MustCompile("\\$[0-9]\\$[a-zA-Z0-9]\\$[a-zA-Z0-9./=]+"), "pwhash"},
		&Keyword{regexp.MustCompile("[a-zA-Z0-9]+::[a-zA-Z0-9]{10}:[a-z0-9]{32}:[a-z0-9-]+"), "pwhash"},
		&Keyword{regexp.MustCompile("[a-zA-Z0-9-_]+:[0-9]+:[a-z0-9]{32}:[a-z0-9]{32}"), "pwhash"},
	}
}
