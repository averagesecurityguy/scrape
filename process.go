package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var cpRegex = regexp.MustCompile("http://.*")
var credRegex = regexp.MustCompile("(?m)^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+:[^ ~/].*$")
var keysRegex = regexp.MustCompile("(?m)BEGIN (RSA|DSA|) PRIVATE KEY.*END")
var awsRegex = regexp.MustCompile("(?i)aws_access_key_id[ =]+(.*)\naws_secret_access_key[ =]+(.*)")

// Find AWS access keys and secrets
func processAWSKeys(contents, url string) {
	keys := awsRegex.FindAllStringSubmatch(contents, -1)

	// No keys found.
	if keys == nil {
		return
	}

	log.Printf("[+] Found AWS keys in: %s", url)

	var formatted []string
	for _, key := range keys {
		formatted = append(formatted, strings.Join(key[1:], ":"))
	}

	save("awskeys.txt", strings.Join(formatted, "\n"))
}

// Look for email addresses and save them to a file.
func processEmails(contents, url string) {
	emails := emailRegex.FindAllString(contents, -1)

	// No emails found.
	if emails == nil {
		return
	}

	// Save the found emails
	log.Printf("[+] Found emails in: %s", url)
	save("emails.txt", strings.Join(emails, "\n"))
}

// Look for credentials in the format of email:password and save them to a file.
func processCredentials(contents, url string) {
	creds := credRegex.FindAllString(contents, -1)

	// No creds found.
	if creds == nil {
		return
	}

	// Save the found creds
	log.Printf("[+] Found credentials in: %s", url)
	save("creds.txt", strings.Join(creds, "\n"))
}

// Look for private keys.
func processPrivKey(contents, url string) {
	keys := keysRegex.FindAllString(contents, -1)

	// No keys found.
	if keys == nil {
		return
	}

	// Save the found keys
	log.Printf("[+] Found private keys in: %s", url)
	save("privkeys.txt", strings.Join(keys, "\n"))
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
func processCopyPaste(contents, title, url string) {
	if !strings.Contains(contents, "Copy & Paste link") {
		return
	}

	log.Printf("[+] Found Copy/Paste link in: %s", url)

	link := cpRegex.FindString(contents)
	if link != "" {
		save("crack_urls.txt", fmt.Sprintf("%s | %s | %s\n", url, title, link))
	}
}

// Save a paste to the data folder with the specified prefix.
func savePaste(prefix string, p *Paste) {
	fname := fmt.Sprintf("%s-%s.paste", prefix, p.Key)

	if p.Expire != 0 && p.Size < conf.maxSize {
		data := fmt.Sprintf("%s\n\n%s", p.Header(), p.Content)
		save(fname, data)
	}
}

func save(fname, data string) {
	path := fmt.Sprintf("%s/%s", conf.dataPath, fname)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[-] Could not open file: %s.", path)
		return
	}

	f.WriteString(data + "\n")
	f.Close()
}

// Process each paste.
func process(p *Paste) {
	// Find and save specific data.
	processEmails(p.Content, p.Url)
	processCredentials(p.Content, p.Url)
	processPrivKey(p.Content, p.Url)
	processCopyPaste(p.Content, p.Title, p.Url)
	processAWSKeys(p.Content, p.Url)

	// Save pastes that match any of our keywords. First match wins. Use these
	// to find interesting data that will eventually be processed with a more
	// specific method.
	for i, _ := range conf.keywords {
		kwd := conf.keywords[i]
		match := kwd.regex.FindString(p.Content)

		if match != "" {
			log.Printf("[+] Found \"%s\" in: %s", kwd.prefix, p.Url)
			savePaste(kwd.prefix, p)
			break
		}
	}
}
