package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Comment struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Priority int    `json:"priority"`
	Text     string `json:"text"`
	Category string `json:"category"`
}

func main() {
	dbURL := fmt.Sprintf("postgres://%s@%s/%s?password=%s&sslmode=disable", "postgres", "localhost", "copypaste", "testing")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err != nil {
			log.Println("Cannot connect to database, retying...")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection is good")

	file, err := ioutil.ReadFile("/home/ben/comments.json")
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}
	var data []Comment
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
	//fmt.Println(data)
	for i := 0; i < len(data); i++ {
		//log.Println("Creating new", resource.Type, "resource:", resource.ID)
		// Prepare the SQL statement
		stmt, err := db.Prepare("INSERT INTO comments (id, label, priority, text, category) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the statement with the item values
		_, err = stmt.Exec(data[i].ID, data[i].Label, data[i].Priority, data[i].Text, data[i].Category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(data[i].Label)
	}
}
