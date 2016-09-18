package metadata

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var metaBucket string = "lincoln-meta"
var root string = fmt.Sprintf("%s/.lincoln", os.Getenv("HOME"))

// Helper for extracting bolt.DB connector
func getDB() *bolt.DB {
	os.Mkdir(root, 0777)
	dbPath := fmt.Sprintf("%s/.lincoln.db", root)
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// ===  Namespaced Buckets  =====================================================
type Namespace []byte

func NS(n string) Namespace {
	return Namespace(n)
}

func (n Namespace) Put(key string, value string) {
	db := getDB()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(n)
		b.Put([]byte(key), []byte(value))
		return nil
	})
}

func (n Namespace) Delete(key string) {
	db := getDB()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(n)
		b.Delete([]byte(key))
		return nil
	})
}

func (n Namespace) Get(key string) string {
	db := getDB()
	defer db.Close()
	var ret []byte

	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(n)
		res := b.Get([]byte(key))
		ret = make([]byte, len(res))
		copy(ret, res)
		return nil
	})

	return string(ret)
}

// ===  Specialized Namespaces  =================================================
func AppNS(n string) Namespace {
	return NS(fmt.Sprintf("app:%v", n))
}

// ===  Metadata Defaults  ======================================================
func PutMeta(key string, value string) {
	NS(metaBucket).Put(key, value)
}

func GetMeta(key string) string {
	return NS(metaBucket).Get(key)
}

func DeleteMeta(key string) {
	NS(metaBucket).Delete(key)
}
