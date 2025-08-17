package main

import (
	"log"
	"net/http"
	"os"

	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
	"github.com/barrynorthern/libretto/services/api/publisher"
	apiserver "github.com/barrynorthern/libretto/services/api/server"
)

func healthMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func main() {
	// Configure port via PORT env var (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	topic := os.Getenv("DIRECTIVE_TOPIC")
	if topic == "" {
		topic = "libretto.dev.directive.issued.v1"
	}
	producer := os.Getenv("PRODUCER")
	if producer == "" {
		producer = "api"
	}

	mux := healthMux()
	pub := publisher.Select()
	// Log which publisher we selected for visibility during manual tests
	switch pub.(type) {
	case publisher.PubSubPublisher:
		log.Printf("publisher=pubsub topic=%s", topic)
	default:
		log.Printf("publisher=nop topic=%s", topic)
	}
	svc := &apiserver.BatonServer{Pub: pub, Topic: topic, Producer: producer}
	mux.Handle(batonv1connect.NewBatonServiceHandler(svc))

	log.Printf("api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
