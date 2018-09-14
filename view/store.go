package main

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type WalkFunc func(key, val string)

type Database struct {
	conn *bolt.DB
}

func NewDatabase(fileName string) (*Database, error) {
	var err error

	db := new(Database)

	for tries := 1; tries < 20; tries += 2 {
		timeout := 1 << uint(tries) * time.Millisecond
		opts := &bolt.Options{Timeout: timeout, ReadOnly: true}

		conn, err := bolt.Open(fileName, 0440, opts)
		if err == nil {
			db.conn = conn
			return db, nil
		}
	}

	return nil, err
}

func (db *Database) WalkBucket(bucket string, walkFn WalkFunc) error {
	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", bucket)
		}

		c := b.Cursor()

		for k, v := c.Seek([]byte("")); k != nil; k, v = c.Next() {
			walkFn(string(k), string(v))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Buckets returns a list of all the buckets in the root of the database.
func (db *Database) Buckets() ([]string, error) {
	var buckets []string

	err := db.conn.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})

		return nil
	})

	if err != nil {
		return buckets, err
	}

	return buckets, nil
}


func (db *Database) Read(bucket, key string) string {
	var val []byte

	db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		val = b.Get([]byte(key))

		return nil
	})

	return string(val)
}

func (db *Database) Close() {
	db.conn.Close()
}
