package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Value struct {
	Bucket string
	Key    string
	Value  string
}

func NewValue(bucket, key string) *Value {
	v := new(Value)

	v.Bucket = bucket
	v.Key = key
	v.Value = db.Read(bucket, key)

	return v
}

type DataSet struct {
	Bucket string
	Batch  map[string]string
	Next   string
}

func NewDataSet(bucket string) *DataSet {
	d := new(DataSet)

	d.Bucket = bucket
	d.Batch = make(map[string]string)

	return d
}

type SearchSet struct {
	Bucket string
	Term   string
	Keys   []string
	Next   string
}

func NewSearchSet(bucket, next, term string) *SearchSet {
	s := new(SearchSet)

	s.Bucket = bucket
	s.Term = term
	s.Next = next

	return s
}

func buckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := db.Buckets()
	if err != nil {
		log.Printf("[-] Could not read buckets: %s\n", err)
		http.Error(w, "could not read buckets", 500)
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/buckets.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.ExecuteTemplate(w, "layout", buckets)
}

func read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/read.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	val := NewValue(vars["bucket"], vars["key"])

	t.ExecuteTemplate(w, "layout", val)
}

func keys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := NewDataSet(vars["bucket"])
	ds.Next = vars["next"]

	err := db.ReadBatch(ds, conf.WebBatchSize)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/keys.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.ExecuteTemplate(w, "layout", ds)
}

func vals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := NewDataSet(vars["bucket"])
	ds.Next = vars["next"]

	err := db.ReadBatch(ds, conf.WebBatchSize)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/vals.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.ExecuteTemplate(w, "layout", ds)
}

func search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ss := NewSearchSet(vars["bucket"], vars["next"], vars["term"])

	for len(ss.Keys) < conf.WebBatchSize {
		temp := NewDataSet(ss.Bucket)
		temp.Next = ss.Next

		err := db.ReadBatch(temp, conf.WebBatchSize)
		if err != nil {
			break
		}

		// Need keys in sorted order to ensure we can set Next correctly.
		var keys []string

		for k := range temp.Batch {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		log.Printf("ss.Keys: %q\n", ss.Keys)
		log.Printf("Keys: %q\n", keys)

		for _, key := range keys {
			if strings.Contains(temp.Batch[key], vars["term"]) {
				ss.Keys = append(ss.Keys, key)
			}
		}

		// If we only have one key in our batch then there are no more keys
		// to find so we are done.
		if len(keys) == 1 {
			ss.Next = ""
			break
		}

		ss.Next = keys[len(keys)-1]
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/search.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t.ExecuteTemplate(w, "layout", ss)
}

func startWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", buckets)
	r.HandleFunc("/keys/{bucket}", keys)
	r.HandleFunc("/keys/{bucket}/{next}", keys)
	r.HandleFunc("/vals/{bucket}", vals)
	r.HandleFunc("/vals/{bucket}/{next}", vals)
	r.HandleFunc("/read/{bucket}/{key}", read)
	r.HandleFunc("/search/{bucket}/{term}", search)
	r.HandleFunc("/search/{bucket}/{term}/{next}", search)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	srv := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         conf.WebServerAddr,
		Handler:      r,
	}

	err := srv.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
	log.Printf("[-] Web server closed: %s\n", err)
}
