package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var header []byte
var footer []byte
var beginList = []byte("<ul>\n")
var endList = []byte("</ul>\n")

func bucketLi(item string) []byte {
	return []byte(fmt.Sprintf("<li><a href=\"/keys/%s\">%s</a></li>\n", item, item))
}

func keyLi(bucket, key string) []byte {
	return []byte(fmt.Sprintf("<li><a href=\"/read/%s/%s\">%s</a></li>\n", bucket, key, key))
}

func valLi(val string) []byte {
	return []byte(fmt.Sprintf("<li>%s</li>\n", val))
}

func pre(data string) []byte {
	return []byte(fmt.Sprintf("<pre>\n%s\n</pre>\n", data))
}

func heading(str string) []byte {
	return []byte(fmt.Sprintf("<h2>%s</h2>\n", str))
}

// buckets  Returns a list of buckets.
func buckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := db.Buckets()
	if err != nil {
		log.Printf("[-] Could not read buckets: %s\n", err)
		http.Error(w, "could not read buckets", 500)
	}

	w.Write(header)
	w.Write(heading("Buckets"))
	w.Write(beginList)
	for b := range buckets {
		w.Write(bucketLi(buckets[b]))
	}
	w.Write(endList)
	w.Write(footer)
}

func read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Write(header)
	w.Write(heading(fmt.Sprintf("%s - %s", vars["bucket"], vars["key"])))
	w.Write(pre(db.Read(vars["bucket"], vars["key"])))
	w.Write(footer)
}

func keys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Write(header)
	w.Write(beginList)
	w.Write(heading(fmt.Sprintf("Keys - %s", vars["bucket"])))
	db.WalkBucket(vars["bucket"], func(key, val string) {
		w.Write(keyLi(vars["bucket"], key))
	})
	w.Write(endList)
	w.Write(footer)
}

func vals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Write(header)
	w.Write(beginList)
	w.Write(heading(fmt.Sprintf("Values - %s", vars["bucket"])))
	db.WalkBucket(vars["bucket"], func(key, val string) {
		w.Write([]byte(valLi(val)))
	})
	w.Write(endList)
	w.Write(footer)
}

func search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Write(header)
	w.Write(beginList)
	w.Write(heading("Search Results"))
	db.WalkBucket(vars["bucket"], func(key, val string) {
		if strings.Contains(val, vars["term"]) {
			w.Write(keyLi(vars["bucket"], key))
		}
	})
	w.Write(endList)
	w.Write(footer)
}

func startWebServer() {
	header, _ = ioutil.ReadFile("web/templates/header.html")
	footer, _ = ioutil.ReadFile("web/templates/footer.html")

	r := mux.NewRouter()
	r.HandleFunc("/", buckets)
	r.HandleFunc("/keys/{bucket}", keys)
	r.HandleFunc("/vals/{bucket}", vals)
	r.HandleFunc("/read/{bucket}/{key}", read)
	r.HandleFunc("/search/{bucket}/{term}", search)

	srv := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         "127.0.0.1:5000",
		Handler:      r,
	}

	err := srv.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
	log.Printf("[-] Web server closed: %s\n", err)
}
