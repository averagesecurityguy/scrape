package main

import (
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type DataSet struct {
	Bucket string
	Batch map[string]string
	Next string
}

func NewDataSet(bucket string) *DataSet {
	d := new(DataSet)

	d.Bucket = bucket
	d.Batch = make(map[string]string)

	return d
}


var header []byte
var footer []byte
var beginList = []byte("<ul>\n")
var endList = []byte("</ul>\n")

func escape(str string) string {
	return html.EscapeString(str)
}

func bucketLi(item string) []byte {
	item = escape(item)

	return []byte(fmt.Sprintf("<li><a href=\"/keys/%s\">%s</a></li>\n", item, item))
}

func pre(data string) []byte {
	data = escape(data)

	return []byte(fmt.Sprintf("<pre>\n%s\n</pre>\n", data))
}

func heading(header string) []byte {
	header = escape(header)

	return []byte(fmt.Sprintf("<h2>%s</h2>\n", header))
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
	w.Write(pre(escape(db.Read(vars["bucket"], vars["key"]))))
	w.Write(footer)
}

func keys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds := NewDataSet(vars["bucket"])
	ds.Next = vars["next"]

	err := db.ReadBatch(ds, 100)
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

	err := db.ReadBatch(ds, 100)
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

	w.Write(header)
	w.Write(heading(fmt.Sprintf("Searching %s for %s", vars["bucket"], vars["term"])))
	w.Write(beginList)
	err := db.WalkBucket(vars["bucket"], func(key, val string) {
		if strings.Contains(val, vars["term"]) {
			w.Write(keyLi(vars["bucket"], key))
		}
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(endList)
	w.Write(footer)
}

func startWebServer() {
	header, _ = ioutil.ReadFile("web/templates/header.html")
	footer, _ = ioutil.ReadFile("web/templates/footer.html")

	r := mux.NewRouter()
	r.HandleFunc("/", buckets)
	r.HandleFunc("/keys/{bucket}", keys)
	r.HandleFunc("/keys/{bucket}/{next}", keys)
	r.HandleFunc("/vals/{bucket}", vals)
	r.HandleFunc("/read/{bucket}/{key}", read)
	r.HandleFunc("/search/{bucket}/{term}", search)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
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
