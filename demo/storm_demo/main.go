package main

import (
	"fmt"
	"log"
	"time"

	"github.com/asdine/storm/v3"
)

type WordEntry struct {
	Word string `storm:"id"`

	SimpleDefinition    string
	DetailedExplanation string

	Proficiency float32 `storm:"index"`

	NextReviewTime int64 `storm:"index"`
}

type User struct {
	ThePrimaryKey string `storm:"id"`     // primary key
	Group         string `storm:"index"`  // this field will be indexed
	Email         string `storm:"unique"` // this field will be indexed with a unique constraint
	Name          string // this field will not be indexed
}

// func main() {
// 	db, err := storm.Open("user.db")
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	defer db.Close()
//
// 	fmt.Println("exec open command ok ")
//
// 	err = db.Init(&User{})
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
//
// 	fmt.Println("exec init ok")
//
// 	Zhang := User{
// 		ThePrimaryKey: "zhangsan",
// 		Group:         "zhangsan",
// 		Email:         "djoafjsdif",
// 		Name:          "fjoadijfo",
// 	}
//
// 	err = db.Save(&Zhang)
//
// 	if err == storm.ErrAlreadyExists {
// 		fmt.Println(Zhang, " has already exist")
// 	}
//
// 	var user User
//
// 	err = db.One("ThePrimaryKey", "zhangsan", &user)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	fmt.Println("user => ", user)
// }

func check(db *storm.DB) {
	var word WordEntry
	err := db.One("Word", "hello", &word)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("word => ", word)
}

func main() {
	db, err := storm.Open("words.db")
	if err != nil {
		fmt.Println("open file")
		log.Fatal(err)
	}

	check(db)

	fmt.Println("exec the open OK")
	// Init the db by struct
	// err = db.Init(&WordEntry{})
	// if err != nil {
	// 	fmt.Println("init db")
	// 	log.Fatal(err)
	// }

	fmt.Println("exec init ok")

	defer db.Close()

	hello := WordEntry{
		Word:                "hello",
		SimpleDefinition:    "nihao",
		DetailedExplanation: "jingchang yongyu dazhaohu deng",
		Proficiency:         0.1,
		NextReviewTime:      time.Now().Unix(),
	}

	err = db.Save(&hello)
	if err == nil {
		fmt.Println("save hello success")
	}
	if err != nil {
		log.Fatal(err)
	}

	err = db.Save(&hello)

	if err == storm.ErrAlreadyExists {
		fmt.Println(err)
	}

	err = db.Update(&WordEntry{
		Word:             "hello",
		SimpleDefinition: "NihaoA",
		Proficiency:      0.8,
	})
	if err != nil {
		log.Fatal(err)
	}
}
