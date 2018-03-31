package main

import (
	"fmt"
	"regexp"
	"strings"
)

var reCreds = regexp.MustCompile("(?m)^([a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+):([^ ~/$].*$)")
var reEmail = regexp.MustCompile("[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]+")
var rePrivKey = regexp.MustCompile("(?s)BEGIN (RSA|DSA|) PRIVATE KEY.*END (RSA|DSA|) PRIVATE KEY")
var reAwsKey = regexp.MustCompile("(?is).*(AKIA[A-Z0-9]{16}).*([A-Za-z0-9+/]{40})")

// Find AWS access keys and secrets
func processAWSKeys(contents, key string) bool {
	awsKeys := reAwsKey.FindAllStringSubmatch(contents, -1)

	// No keys found.
	if awsKeys == nil {
		return false
	}

	for _, awsKey := range awsKeys {
		conf.ds.Write("awskeys", strings.Join(awsKey[1:], ":"), []byte(key))
	}

	return true
}

// Look for email addresses and save them to a file.
func processEmails(contents, key string) bool {
	emails := reEmail.FindAllString(contents, -1)

	// No emails found.
	if emails == nil {
		return false
	}

	for _, email := range emails {
		email = cleanEmail(email)

		if email == "" {
			continue
		}

		conf.ds.Write("emails", email, []byte(key))
	}

	return true
}

// Look for credentials in the format of email:password and save them to a file.
func processCredentials(contents, key string) bool {
	creds := reCreds.FindAllString(contents, -1)

	// No creds found.
	if creds == nil {
		return false
	}

	for _, cred := range creds {
		conf.ds.Write("creds", cred, []byte(key))
	}

	return true
}

// Look for private keys.
func processPrivKey(contents, key string) bool {
	privKeys := rePrivKey.FindAllString(contents, -1)

	// No keys found.
	if privKeys == nil {
		return false
	}

	for _, privKey := range privKeys {
		conf.ds.Write("privkeys", privKey, []byte(key))
	}

	return true
}

func processContent(key, content string) {
	// Find and save specific data.
	switch {
	case processCredentials(content, key):
		conf.ds.Write("pastes", key, []byte(content))
	case processEmails(content, key):
		conf.ds.Write("pastes", key, []byte(content))
	case processPrivKey(content, key):
		conf.ds.Write("pastes", key, []byte(content))
	case processAWSKeys(content, key):
		conf.ds.Write("pastes", key, []byte(content))
	default:
	}

	// Save pastes that match any of our regular expressions. Use these to find
	// interesting data that will eventually be processed with a more specific
	// method.
	save := false
	for i, _ := range conf.Regexes {
		r := conf.Regexes[i]
		rKey := fmt.Sprintf("%s-%s", r.Prefix, key)
		match := r.compiled.FindString(content)

		if match != "" {
			save = true
			conf.ds.Write("regexes", rKey, nil)
		}
	}

	if save {
		conf.ds.Write("pastes", key, []byte(content))
	}

	// Save pastes that match any of our keywords. Use these to find interesting
	// data that will eventually be processed with a more specific method.
	save = false
	for i, _ := range conf.Keywords {
		kwd := conf.Keywords[i]
		kwdKey := fmt.Sprintf("%s-%s", kwd.Prefix, key)

		if strings.Contains(strings.ToLower(content), kwd.Keyword) {
			save = true
			conf.ds.Write("keywords", kwdKey, nil)
		}
	}

	if save {
		conf.ds.Write("pastes", key, []byte(content))
	}
}

// Remove common false positives in email addresses.
func cleanEmail(email string) string {
	email = strings.ToLower(email)

	switch {
	case strings.HasSuffix(email, "2x.png"):
		return ""
	case strings.HasSuffix(email, ".so"):
		return ""
	default:
		return email
	}
}
