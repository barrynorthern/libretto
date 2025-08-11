package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for Pub/Sub push message handling
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("received"))
}

func main() {
	http.HandleFunc("/", handler)
	addr := ":8080"
	log.Printf("plotweaver listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

