package main

import (
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/gorilla/mux"
)

// buckets  Returns a list of buckets.
func buckets(w http.ResponseWriter, r *http.Request) {
    buckets, err := db.Buckets()
    if err != nil {
        log.Printf("[-] Could not read buckets: %s\n", err)
    	http.Error(w, "could not read buckets", 500)
    }

    for b := range buckets {
    	w.Write([]byte(buckets[b]))
        w.Write([]byte("\n"))
    }
}

func read(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    w.Write([]byte(db.Read(vars["bucket"], vars["key"])))
}

func keys(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

	db.WalkBucket(vars["bucket"], func(key, val string) {
        w.Write([]byte(key))
        w.Write([]byte("\n"))
	})
}

func vals(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

	db.WalkBucket(vars["bucket"], func(key, val string) {
		w.Write([]byte(val))
        w.Write([]byte("\n"))
	})
}

func search(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

	db.WalkBucket(vars["bucket"], func(key, val string) {
		if strings.Contains(val, vars["term"]) {
			w.Write([]byte(val))
            w.Write([]byte("\n"))
		}
	})
}

func startWebServer() {
    // tmpls = template.New("template")
    // tmpls = tmpls.Funcs(template.FuncMap{"buildUrl": buildUrl})
    // tmpls = template.Must(tmpls.ParseFiles(tFiles...))

    r := mux.NewRouter()
    r.HandleFunc("/buckets/", buckets)
    r.HandleFunc("/keys/{bucket}", keys)
    r.HandleFunc("/vals/{bucket}", vals)
    r.HandleFunc("/read/{bucket}/{key}", read)
    r.HandleFunc("/search/{bucket}/{term}", search)

    srv := &http.Server{
        ReadTimeout: 15 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout: 120 * time.Second,
        Addr: "127.0.0.1:5000",
        Handler: r,
    }

    srv.ListenAndServe()
}
