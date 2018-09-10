package main

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type Database struct {
	conn *bolt.DB
}

func newDatabase(fileName string) *Database {
	db := new(Database)

	for tries := 1; tries < 20; tries += 2 {
		timeout := 1 << uint(tries) * time.Millisecond

		conn, err := bolt.Open(fileName, 0640, &bolt.Options{Timeout: timeout})
		if err == nil {
			db.conn = conn
			return db
		}

		log.Printf("[-] Database locked waiting...\n")
	}

	log.Printf("[-] Could not connect to the database: %s\n", fileName)

	return nil
}

func (db *Database) CreateBucket(bucket string) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return nil
	})
}

func (db *Database) Write(bucket, key string, value []byte) {
	err := db.conn.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put([]byte(key), []byte(value))
	})

	if err != nil {
		log.Printf("[-] Could not write key: %s\n", err)
	}
}

func (db *Database) Read(bucket, key string) []byte {
	var val []byte

	db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		val = b.Get([]byte(key))

		return nil
	})

	return val
}

func (db *Database) Delete(bucket, key string) {
	err := db.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(key))
	})

	if err != nil {
		log.Printf("[-] Could not delete key: %s\n", err)
	}
}
