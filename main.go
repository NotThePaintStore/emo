package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./test.json")
}

func main() {
	http.Handle("/", http.RedirectHandler("https://emo.bmoore.xyz/copypaste/", http.StatusSeeOther))
	fs := http.FileServer(http.Dir("./copypaste/"))
	http.Handle("/copypaste/", http.StripPrefix("/copypaste", fs))
	http.HandleFunc("/comments.json", serveJSON)
	fmt.Println("Server active on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
