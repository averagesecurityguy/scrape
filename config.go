package main

import (
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
	dataPath string
	maxSize  int           // Do not save files larger than this.
	maxTime  time.Duration // Max time to store previously downloaded keys.
	sleep    time.Duration // Time to wait between each run.
}

func newConfig() Config {
	var c Config

	c.keys = make(map[string]time.Time)
	c.maxSize = 1024 * 1024
	c.maxTime = 3600 * time.Second
	c.sleep = 60 * time.Second
	c.dataPath = "data"

	// Build our includes list.
	c.keywords = loadKeywords()

	return c
}

func loadKeywords() []*Keyword {
	return []*Keyword{
		&Keyword{regexp.MustCompile("(?i)BEGIN PRIVATE KEY"), "privkey"},
		&Keyword{regexp.MustCompile("(?i)BEGIN DSA PRIVATE KEY"), "privkey"},
		&Keyword{regexp.MustCompile("(?i)BEGIN RSA PRIVATE KEY"), "privkey"},
		&Keyword{regexp.MustCompile("(?i)FULLZ"), "carder"},
		&Keyword{regexp.MustCompile("(?i)aws_secret_access_key"), "awskey"},
		&Keyword{regexp.MustCompile("(?i)`password`"), "sqlpass"},
	}
}
