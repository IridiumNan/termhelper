package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/boltdb/bolt"
)

const (
	dbName              = "words.db"
	SimpleDefBucketName = "def"
	detailedExplainName = "det"
	proficiencyName     = "pro"
)

type WordEntry struct {
	Word                string
	SimpleDefinition    string
	DetailedExplanation string

	Proficienty float32
}

func BytesFloat32(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)

	return float
}

func createBucket() (*bolt.DB, error) {
	db, err := bolt.Open(dbName, 0o600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(SimpleDefBucketName))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(detailedExplainName))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(proficiencyName))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertNewWordEntry(db *bolt.DB, word WordEntry) error {
	err := db.Update(func(tx *bolt.Tx) error {
		defBucket := tx.Bucket([]byte(SimpleDefBucketName))

		err := defBucket.Put([]byte(word.Word), []byte(word.SimpleDefinition))
		if err != nil {
			return err
		}

		detBucket := tx.Bucket([]byte(detailedExplainName))

		err = detBucket.Put([]byte(word.Word), []byte(word.DetailedExplanation))
		if err != nil {
			return err
		}

		proBucket := tx.Bucket([]byte(proficiencyName))

		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, math.Float32bits(word.Proficienty))

		err = proBucket.Put([]byte(word.Word), bytes)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func GetWordEntry(db *bolt.DB, wordStr string) (*WordEntry, error) {
	var word *WordEntry
	err := db.View(func(tx *bolt.Tx) error {
		defBucket := tx.Bucket([]byte(SimpleDefBucketName))

		wordKey := []byte(wordStr)

		simpleDef := defBucket.Get(wordKey)
		if simpleDef == nil {
			return fmt.Errorf("key %s not found", wordStr)
		}

		detBucket := tx.Bucket([]byte(detailedExplainName))

		detailExp := detBucket.Get(wordKey)

		if detailExp == nil {
			return fmt.Errorf("key %s not found", wordStr)
		}

		proBucket := tx.Bucket([]byte(proficiencyName))

		proBytes := proBucket.Get(wordKey)

		if proBytes == nil {
			return fmt.Errorf("key %s not found", wordStr)
		}

		pro := BytesFloat32(proBytes)

		word = &WordEntry{
			Word:                wordStr,
			SimpleDefinition:    string(simpleDef),
			DetailedExplanation: string(detailExp),
			Proficienty:         pro,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return word, nil
}

func main() {
	db, err := createBucket()
	if err != nil {
		log.Fatal(err)
	}

	// word := WordEntry{
	// 	Word:                "hello",
	// 	SimpleDefinition:    "你好",
	// 	DetailedExplanation: "一般用于打招呼， 常见用法有hello world",
	// 	Proficienty:         0.5,
	// }
	//
	// err = insertNewWordEntry(db, word)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	hello, err := GetWordEntry(db, "hello")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hello entry ->", hello)

	no, err := GetWordEntry(db, "no")
	if err != nil {
		log.Fatal(err)
		db.Close()
	}

	fmt.Println("no entry -> ", no)

	defer db.Close()
}
