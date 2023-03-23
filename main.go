package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

var urlMap = make(map[string]string)

func shortenURL(longURL string) string {
	hash := md5.Sum([]byte(longURL))
	shortURL := hex.EncodeToString(hash[:])[:8]
	urlMap[shortURL] = longURL
	return shortURL
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, ok := urlMap[shortURL]

	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	longURL := r.URL.Query().Get("url")
	if longURL == "" {
		http.Error(w, "url parameter is missing", http.StatusBadRequest)
		return
	}

	testURL := "l2.fm"

	shortURL := shortenURL(longURL)
	fmt.Fprintf(w, "http://%s/%s\n", r.Host, shortURL)
	fmt.Fprintf(w, "%s/%s\n", testURL, shortURL)
}

func main() {
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
