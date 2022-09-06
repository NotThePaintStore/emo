package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.RedirectHandler("https://emo.bmoore.xyz/copypaste/", http.StatusSeeOther))
	fs := http.FileServer(http.Dir("./copypaste/"))
	http.Handle("/copypaste/", http.StripPrefix("/copypaste", fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}