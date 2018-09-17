package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type WalkFunc func(key, val string)

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

func (db *Database) ReadBatch(ds *DataSet, count int) error {
	 return db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ds.Bucket))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", ds.Bucket)
		}

		c := b.Cursor()

		for k, v := c.Seek([]byte(ds.Next)); k != nil && len(ds.Batch) <= count; k, v = c.Next() {
			ds.Batch[string(k)] = string(v)
			ds.Next = string(k)
		}

		if len(ds.Batch) != count {
			ds.Next = ""
		}

		return nil
	})
}

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

func (db *Database) Delete(bucket, key string) {
	err := db.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(key))
	})

	if err != nil {
		log.Printf("[-] Could not delete key: %s\n", err)
	}
}
