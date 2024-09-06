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

func allHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	comments, err := GetAllComments(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func activeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	comments, err := GetAllActiveComments(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func updateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		var comment Comment
		comment.ID = r.FormValue("commentID")
		comment.Label = r.FormValue("newLabel")
		comment.Priority, _ = strconv.Atoi(r.FormValue("priority"))
		comment.Category = r.FormValue("category")
		comment.Text = r.FormValue("text")
		activestat := r.FormValue("active")
		if activestat == "yes" {
			comment.Active = true
		} else {
			comment.Active = false
		}

		err = UpdateComment(db, comment)
		if err != nil {
			http.Error(w, "Error updating comment", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s has been updated.", comment.Label)
	}
}

func singleHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	queryParams := r.URL.Query()
	commentID := queryParams.Get("id")
	comment, err := GetComment(db, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetComment(db *sql.DB, commentID string) (*Comment, error) {
	log.Println("Getting details for", commentID)
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT (id || '|' || label || '|' || priority || '|' || text || '|' || category || '|' || active) FROM comments WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var s string
	var res Comment

	// Execute the statement with the comment ID
	results := stmt.QueryRow(commentID)
	err = results.Scan(&s)
	if err != nil {
		return nil, err
	}

	// Format data and put it in a Comment struct
	sl := strings.Split(strings.Trim(s, "()"), "|")
	res.Priority, _ = strconv.Atoi(sl[2])
	res.Label = strings.Trim(sl[1], `"`)
	res.ID = sl[0]
	res.Text = strings.Trim(sl[3], `"`)
	res.Category = strings.Trim(sl[4], `"`)
	res.Active, _ = strconv.ParseBool(sl[5])

	return &res, nil
}

func UpdateComment(db *sql.DB, comment Comment) error {
	log.Println("Updating", comment.ID)
	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE comments SET label = $1, priority = $2, text = $3, active = $4, category = $5 WHERE id = $6")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the item values
	_, err = stmt.Exec(comment.Label, comment.Priority, comment.Text, comment.Active, comment.Category, comment.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllActiveComments(db *sql.DB) ([]Comment, error) {
	log.Println("Getting active comments")
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT (id || '|' || label || '|' || priority || '|' || text || '|' || category || '|' || active) FROM comments WHERE active IS true ORDER BY priority")
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
			log.Println(err)
			continue
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
		res.ID = sl[0]
		res.Text = strings.Trim(sl[3], `"`)
		res.Category = strings.Trim(sl[4], `"`)
		res.Active, _ = strconv.ParseBool(sl[5])
		comments = append(comments, res)
	}
	return comments, nil
}

func GetAllComments(db *sql.DB) ([]Comment, error) {
	log.Println("Getting all comments")
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT (id || '|' || label || '|' || priority || '|' || text || '|' || category || '|' || active) FROM comments ORDER BY label")
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
			log.Println(err)
			continue
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
		res.ID = sl[0]
		res.Text = strings.Trim(sl[3], `"`)
		res.Category = strings.Trim(sl[4], `"`)
		res.Active, _ = strconv.ParseBool(sl[5])
		comments = append(comments, res)
	}
	return comments, nil
}

func main() {
	dbURL := fmt.Sprintf("postgres://%s@%s/%s?password=%s&sslmode=disable", "postgres", "copypaste-db", "copypaste", "testing")
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

	http.Handle("/home", fs)
	http.HandleFunc("/active", func(w http.ResponseWriter, r *http.Request) { activeHandler(w, r, db) })
	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) { allHandler(w, r, db) })
	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) { singleHandler(w, r, db) })
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) { updateHandler(w, r, db) })
	http.Handle("/", fs)

	fmt.Println("Server active on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
