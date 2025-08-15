package main

import (
	"log"
	"net/http"
	"os"

	"github.com/barrynorthern/libretto/gen/go/libretto/graph/v1/graphv1connect"
	gwserver "github.com/barrynorthern/libretto/services/graphwrite/server"
)

func main() {
	mux := http.NewServeMux()
	store := &gwserver.InMemoryStore{}
	svc := &gwserver.GraphWriteServer{Store: store}
	mux.Handle(graphv1connect.NewGraphWriteServiceHandler(svc))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	addr := ":" + port
	log.Printf("graphwrite listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
