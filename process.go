package main

import (
	"fmt"
	"regexp"
	"strings"
	"os"
)

// Look for credentials in the format of email:password and save them to a file.
func processCredentials(contents string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9-+_.]+@[a-zA-Z0-9.-]+:.*")
	creds := re.FindAllString(contents, -1)

	// No creds found.
	if creds == nil {
		return false
	}

	// Save the found creds
	f, err := os.OpenFile("data/creds.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
        fmt.Println("Could not open creds file.")
		return true
    }

	for _, cred := range creds {
		f.WriteString(fmt.Sprintf("%s\n", cred))
	}

	f.Close()

	return true
}

// Found a lot of files with the format:
//
//
// ********************
// Tengo Problemas Para Entrar A Skype
// http://tinyurl.com/y7ghsneu
// (Copy & Paste link)
// ********************
//
// ...
// Keywords
//
// Example: https://pastebin.com/GP7Gx41u
// This method extracts those URLs for later analysis.
func processCopyPaste(purl, title, contents string) {
	re := regexp.MustCompile("http://.*")
	url := re.FindString(contents)

	if url != "" {
		// Save the found url
		f, err := os.OpenFile("data/crack_urls.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
	        fmt.Println("Could not open crack urls file.")
			return
	    }

		f.WriteString(fmt.Sprintf("%s|%s|%s\n", purl, title, url))

		f.Close()
	}
}

// Save a paste to the data folder with the specified prefix.
func save(prefix, key, data string) {
	fname := fmt.Sprintf("data/%s-%s.paste", prefix, key)

	fd, err := os.Create(fname)
	if err != nil {
		fmt.Printf("Could not create file: %s\n", err.Error())
		return
	}

	defer fd.Close()

	fd.WriteString(data)
	fd.Close()
}

// Process each paste.
func process(p Paste) {
	if processCredentials(p.Content) {
		save("creds", p.Key, p.String())
		return
	}

	if strings.Contains(p.Content, "Copy & Paste link") {
	     processCopyPaste(p.Url, p.Title[:25], p.Content)
		 save("cp", p.Key, p.String())
		 return
    }

	// Save pastes that have any of our keywords. First match wins.
	for i, _ := range conf.keywords {
		kwd := conf.keywords[i]

		if strings.Contains(p.Content, kwd.word) {
			save(kwd.prefix, p.Key, p.String())
			break
		}
	}
}
