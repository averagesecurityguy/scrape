package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func scrapeFiles(c chan<- *ProcessItem) {
	if conf.LocalPath == "" {
		return
	}

	log.Println("[+] Checking for local pastes.")

	files, err := ioutil.ReadDir(conf.LocalPath)
	if err != nil {
		log.Printf("[-] Error reading %s: %s\n", conf.LocalPath, err)
		return
	}

	// Process files in batches
	for _, file := range files[:conf.FileBatchSize] {
		if file.IsDir() {
			log.Printf("[+] Skipping directory %s\n", conf.LocalPath)
			continue
		}

		path := filepath.Join(conf.LocalPath, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("[-] Could not read file %s\n", path)
			continue
		}

		item := &ProcessItem{Source: "Local", Key: file.Name(), Content: string(content)}
		c <- item

		os.Remove(path)
	}
}
