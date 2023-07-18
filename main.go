package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Comment struct {
	ID       string
	Label    string
	Priority int
	Text     string
	Category string
	Active   bool
}

func activeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	comments, err := GetAllActiveComments(db)
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetAllActiveComments(db *sql.DB) ([]Comment, error) {
	log.Println("Getting all comments")
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT (id || '|' || label || '|' || priority || '|' || text || '|' || category || '|' || active) FROM comments WHERE active IS true")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement with the item values
	results, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer results.Close()

	// Put results into a slice
	rows := make([]string, 0)
	for results.Next() {
		var row string
		if err := results.Scan(&row); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	// make a slice of Comments to return
	var comments []Comment
	for i := 0; i < len(rows); i++ {
		var res Comment
		// Format data and put it in a Comment struct
		sl := strings.Split(strings.Trim(rows[i], "()"), "|")
		res.Priority, _ = strconv.Atoi(sl[2])
		res.Label = strings.Trim(sl[1], `"`)
		res.ID = strings.Trim(sl[0], `"`)
		res.Text = strings.Trim(sl[3], `"`)
		res.Category = strings.Trim(sl[4], `"`)
		res.Active, _ = strconv.ParseBool(sl[5])
		comments = append(comments, res)
	}
	return comments, nil
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

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/", fs)
	http.HandleFunc("/active", func(w http.ResponseWriter, r *http.Request) { activeHandler(w, r, db) })
	//http.HandleFunc("/comments", serveJSON)
	fmt.Println("Server active on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
