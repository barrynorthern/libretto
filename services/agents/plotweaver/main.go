package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	// Configure port via PORT env var (default 8081 to avoid clashing with API)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	log.Printf("plotweaver listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
