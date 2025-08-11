package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	addr := ":8080"
	log.Printf("api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

