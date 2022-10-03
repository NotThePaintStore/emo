package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./copypaste/comments.json")
}

func main() {
	fs := http.FileServer(http.Dir("./copypaste/"))
	http.Handle("/", fs)
	http.HandleFunc("/comments.json", serveJSON)
	fmt.Println("Server active on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
