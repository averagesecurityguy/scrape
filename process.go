package main

import (
	"log"
	"regexp"
	"strings"
)

var reCreds = regexp.MustCompile("(?m)^([a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+):([^ ~/$].*$)")
var reEmail = regexp.MustCompile("[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]+")
var rePrivKey = regexp.MustCompile("(?s)BEGIN (RSA|DSA|) PRIVATE KEY.*END (RSA|DSA|) PRIVATE KEY")
var reAwsKey = regexp.MustCompile("(?is).*(AKIA[A-Z0-9]{16}).*([A-Za-z0-9+/]{40})")
var ds *KVStore

// Find AWS access keys and secrets
func processAWSKeys(contents string) bool {
	keys := reAwsKey.FindAllStringSubmatch(contents, -1)

	// No keys found.
	if keys == nil {
		return false
	}

	for _, key := range keys {
		ds.Put("awskeys", key[1], key[2])
	}

	return true
}

// Look for email addresses and save them to a file.
func processEmails(contents string) bool {
	emails := reEmail.FindAllString(contents, -1)

	// No emails found.
	if emails == nil {
		return false
	}

	for _, email := range emails {
		ds.Put("emails", strings.ToLower(email), "")
	}

	return true
}

// Look for credentials in the format of email:password and save them to a file.
func processCredentials(contents string) bool {
	creds := reCreds.FindAllString(contents, -1)

	// No creds found.
	if creds == nil {
		return false
	}

	for _, cred := range creds {
		ds.Put("creds", cred, "")
	}

	return true
}

// Look for private keys.
func processPrivKey(contents string) bool {
	keys := rePrivKey.FindAllString(contents, -1)

	// No keys found.
	if keys == nil {
		return false
	}

	for _, key := range keys {
		ds.Put("privkeys", key, "")
	}

	return true
}

// Process each paste.
func process(p *Paste) {
	ds, err := NewKVStore(conf.dbPath)
	if err != nil {
		log.Printf("[-] Cannot open database. Skipping processing.")
		return
	}

	// Find and save specific data.
	if processEmails(p.Content) || processCredentials(p.Content) ||
		processPrivKey(p.Content) || processAWSKeys(p.Content) {
		ds.Put("rawpastes", p.Key, p)
	}

	// Save pastes that match any of our keywords. First match wins. Use these
	// to find interesting data that will eventually be processed with a more
	// specific method.
	for i, _ := range conf.keywords {
		kwd := conf.keywords[i]
		match := kwd.regex.FindString(p.Content)

		if match != "" {
			ds.Put("keywords", kwd.prefix, p)
			break
		}
	}
}
