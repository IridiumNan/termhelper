package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type WordEntry struct {
	Word          string  `json:"word"`
	Definition    string  `json:"definition"`
	Proficiency   float32 `json:"proficiency"`
	NextReviewDue int64   `json:"next_review_due"`
	CreateAt      int64   `json:"created_at"`
}

func main() {
	// quit if cannot get lock in 1 second
	db, err := bolt.Open("my.db", 0o600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("words"))

		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	entry := WordEntry{
		Word:        "compile",
		Definition:  "编译",
		Proficiency: 0.15,

		NextReviewDue: time.Now().Add(24 * time.Hour).Unix(),

		CreateAt: time.Now().Unix(),
	}

	value, _ := json.Marshal(entry)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("words"))

		return b.Put([]byte("compile"), value)
	})
	if err != nil {
		log.Fatal(err)
	}
}
