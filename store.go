package main

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func getDBConn() *bolt.DB {
	for tries := 1; tries < 20; tries += 2 {
		timeout := 1 << uint(tries) * time.Millisecond

		db, err := bolt.Open(conf.DbFile, 0640, &bolt.Options{Timeout: timeout})
		if err == nil {
			return db
		}

		log.Printf("[-] Database locked waiting...\n")
	}

	return nil
}

func initDB() {
	db := getDBConn()

	createBucket(db, "pastes")

	for _, kw := range conf.Keywords {
		createBucket(db, kw.Bucket)
	}

	for _, re := range conf.Regexes {
		createBucket(db, re.Bucket)
	}

	db.Close()
}

func createBucket(db *bolt.DB, bucket string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return nil
	})
}

func writeDB(db *bolt.DB, bucket, key string, value []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put([]byte(key), []byte(value))
	})

	if err != nil {
		log.Printf("[-] Could not write key: %s\n", err)
	}
}

func readDB(db *bolt.DB, bucket, key string) []byte {
	var val []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		val = b.Get([]byte(key))

		return nil
	})

	return val
}

func deleteDB(db *bolt.DB, bucket, key string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(key))
	})

	if err != nil {
		log.Printf("[-] Could not delete key: %s\n", err)
	}
}
