package main

import (
	"fmt"
	"os"
	"strings"
)

func help() {
	u := `Usage:
    view filename action [arguments]

Actions:
    buckets                       Get a list of buckets.
    read <bucketname> <key>       Get the value of the key in the bucket.
    keys <bucketname>             Get a list of keys in a bucket.
    vals <bucketname>             Get a list of values in a bucket.
    search <bucketname> <string>  Get a list of keys from the bucket where the
                                  value contains the given string.
	`
	fmt.Println(u)
	os.Exit(0)
}

// buckets  Returns a list of buckets.
func buckets(db *Database) {
	buckets, err := db.Buckets()
	if err != nil {
		fmt.Printf("[-] Could not read buckets: %s.\n", err)
	}

	for b := range buckets {
		fmt.Println(buckets[b])
	}
}

// read <bucketname> <key>  Returns the value of the key in the bucket.
func read(db *Database, args []string) {
	switch len(args) {
	case 2:
		fmt.Println(db.Read(args[0], args[1]))
	default:
		help()
	}
}

// keys <bucketname>  Returns all keys in a bucket.
func keys(db *Database, args []string) {
	switch len(args) {
	case 1:
		db.WalkBucket(args[0], func(key, val string) {
			fmt.Println(key)
		})
	default:
		help()
	}
}

// vals <bucketname>  Returns all vals in a bucket.
func vals(db *Database, args []string) {
	switch len(args) {
	case 1:
		db.WalkBucket(args[0], func(key, val string) {
			fmt.Println(val)
		})
	default:
		help()
	}
}

// search <bucketname> <string> Returns all keys in a bucket whose value contains the given string.
func search(db *Database, args []string) {
	switch len(args) {
	case 2:
		db.WalkBucket(args[0], func(key, val string) {
			if strings.Contains(val, args[1]) {
				fmt.Println(key)
			}
		})
	default:
		help()
	}
}

func main() {
	if len(os.Args) < 3 {
		help()
	}

	// Open our database file.
	dbfile := os.Args[1]
	db, err := NewDatabase(dbfile)
	if err != nil {
		fmt.Println("[-] Could not open database file:", err)
		os.Exit(0)
	}

	action := os.Args[2]

	switch action {
	case "buckets":
		buckets(db)
	case "read":
		read(db, os.Args[3:])
	case "keys":
		keys(db, os.Args[3:])
	case "vals":
		vals(db, os.Args[3:])
	case "search":
		search(db, os.Args[3:])
	default:
		help()
	}

	db.Close()
}
