package main

import (
	"log"
	"net/http"
	"os"

	"github.com/barrynorthern/libretto/services/agents/plotweaver/publisher"
)

var (
	plotPublisher publisher.Publisher
)

func main() {
	// Publisher selection
	var sel string
	plotPublisher, sel = publisher.Select()
	log.Printf("plotweaver publisher=%s", sel)

	http.HandleFunc("/", handler)
	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	// Configure port via PORT env var (default 8081 to avoid clashing with API)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	log.Printf("plotweaver listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
