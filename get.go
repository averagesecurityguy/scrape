package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func get(url string) []byte {
	return processHTTP(http.Get(url))
}

func getGithub(url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[-] Could not create HTTP request: %s\n", err.Error())
		return []byte("")
	}

	req.Header.Set("Authorization", "token "+conf.GhToken)
	resp, err := client.Do(req)
	return processHTTP(resp, err)
}

func processHTTP(resp *http.Response, err error) []byte {
	if err != nil {
		log.Println("[-] Could not access url.")
		return []byte("")
	}

	if resp.StatusCode != 200 {
		log.Printf("[-] Received HTTP error %d.\n", resp.StatusCode)
		return []byte("")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[-] Could not read response")
		return []byte("")
	}

	return body
}
