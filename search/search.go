package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/asggo/store"
)

func main() {
	var dbfile string
	var bucket string
	var key string
	var query string

    flag.StringVar(&dbfile, "d", "", "Path to Store file.")
	flag.StringVar(&bucket, "b", "", "Search the keys in this bucket.")
	flag.StringVar(&key, "k", "", "Retrieve the value of this key from the bucket.")
	flag.StringVar(&query, "q", "", "List keys that match this query.")

	flag.Parse()

	if dbfile == "" {
		fmt.Println("No database provided.")
		flag.Usage()
	}

	db, err := store.NewStore(dbfile)

	switch {
	case err != nil:
		fmt.Println("Could not open database file:", err)
	case key != "":
		val := db.Read(bucket, key)
		fmt.Printf("%s:%s - %s\n", bucket, key, val)
	default:
		keys, err := db.FindKeys(bucket, query)
		if err != nil {
			fmt.Println("Error retrieving keys:", err)
		} else {
			fmt.Println(strings.Join(keys, "\n"))
		}
	}
}
