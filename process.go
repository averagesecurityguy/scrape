package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
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
		fmt.Println("[-] Could not open creds file.")
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
			fmt.Println("[-] Could not open crack urls file.")
			return
		}

		f.WriteString(fmt.Sprintf("%s|%s|%s\n", purl, title, url))

		f.Close()
	}
}

// Save a paste to the data folder with the specified prefix.
func save(prefix string, p *Paste) {
	fname := fmt.Sprintf("data/%s-%s.paste", prefix, p.Key)

	fd, err := os.Create(fname)
	if err != nil {
		fmt.Printf("[-] Could not create file: %s\n", err.Error())
		return
	}

	defer fd.Close()

	// Do not save pastes that do not expire. We can look them up later.
	if p.Expire == 0 {
		fd.WriteString(p.Header())
	} else {
		fd.WriteString(p.Header())

		// Do not save large pastes. Don't want to fill up disk.
		if p.Size < conf.maxSize {
			fd.WriteString(p.Content)
		}
	}
	fd.Close()
}

// Process each paste.
func process(p *Paste) {
	if processCredentials(p.Content) {
		fmt.Printf("[+] Found credentials in: %s\n", p.Url)
		save("creds", p)
		return
	}

	if strings.Contains(p.Content, "Copy & Paste link") {
		fmt.Printf("[+] Found Copy/Paste link in: %s\n", p.Url)
		processCopyPaste(p.Url, p.Title[:25], p.Content)
		save("cp", p)
		return
	}

	// Save pastes that have any of our keywords. First match wins.
	for i, _ := range conf.keywords {
		kwd := conf.keywords[i]

		if strings.Contains(p.Content, kwd.word) {
			fmt.Printf("[+] Found \"%s\" in: %s\n", kwd.word, p.Url)

			save(kwd.prefix, p)
			break
		}
	}
}
