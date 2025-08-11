package main

import (
	"log"
	"net/http"
)

func healthHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func main() {
	addr := ":8080"
	log.Printf("api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, healthHandler()))
}

