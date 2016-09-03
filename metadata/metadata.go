package metadata

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var metaBucket []byte = []byte("lincoln-meta")

func getDB() *bolt.DB {
	path := fmt.Sprintf("%s/.lincoln.db", os.Getenv("HOME"))
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func PutMeta(key string, value string) {
	db := getDB()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(metaBucket)
		b.Put([]byte(key), []byte(value))
		return nil
	})
}

func GetMeta(key string) string {
	db := getDB()
	defer db.Close()
	var ret []byte

	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(metaBucket)
		res := b.Get([]byte(key))
		ret = make([]byte, len(res))
		copy(ret, res)
		return nil
	})

	return string(ret)
}
