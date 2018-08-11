package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Path    string
	Key     string
	Content string
}

func (f *File) Read() {
	content, err := ioutil.ReadFile(f.Path)
	if err != nil {
		log.Printf("[-] Could not read file %s\n", f.Path)
		f.Content = ""
	} else {
		f.Content = string(content)
	}
}

func (f *File) Process() {
	processContent(f.Key, f.Content)
}

func (f *File) Delete() {
	os.Remove(f.Path)
}

func scrapeFiles() {
	if conf.LocalPath == "" {
		return
	}

	var files []*File

	log.Println("[+] Checking for local pastes.")

	filepath.Walk(conf.LocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[-] Error reading %s: %s\n", conf.LocalPath, err.Error())
			return nil
		}

		if info.IsDir() {
			log.Printf("[+] Skipping directory %s\n", path)
			return nil
		}

		files = append(files, &File{Path: path, Key: filepath.Base(path)})

		return nil
	})

	for i := range files {
		f := files[i]
		f.Read()
		f.Process()
		f.Delete()
	}
}
