package main

import (
	"log"
	"net/http"
	"os"

	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
	"github.com/barrynorthern/libretto/internal/app"
	gwpkg "github.com/barrynorthern/libretto/internal/graphwrite"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	mux := http.NewServeMux()

	// Wire orchestrated Baton service
	orchestrator := app.NewOrchestrator()
	mux.Handle(batonv1connect.NewBatonServiceHandler(orchestrator))
	// JSON endpoints for UI
	app.RegisterHTTP(mux, gwpkg.NewInMemory())

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("libretto (monolith) listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
