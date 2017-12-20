package main

import (
	"time"
)

type Keyword struct {
	word   string
	prefix string
}

type Config struct {
	keys     map[string]time.Time
	keywords []*Keyword
	maxSize  int  // Do not save files larger than this.
	maxTime  time.Duration // Max time to store previously downloaded keys.
	sleep    time.Duration // Time to wait between each run.
}

func newConfig() Config {
	var c Config

	c.keys = make(map[string]time.Time)
	c.maxSize = 1024 * 1024
	c.maxTime = 3600 * time.Second
	c.sleep = 60 * time.Second

	// Build our includes list.
	c.keywords = append(c.keywords, &Keyword{"BEGIN PRIVATE KEY", "privkey"})
	c.keywords = append(c.keywords, &Keyword{"BEGIN RSA PRIVATE KEY", "privkey"})
	c.keywords = append(c.keywords, &Keyword{"BEGIN DSA PRIVATE KEY", "privkey"})
	c.keywords = append(c.keywords, &Keyword{"FULLZ", "carders"})

	return c
}
