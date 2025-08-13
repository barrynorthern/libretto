package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	addr := ":8080"
	log.Printf("plotweaver listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
