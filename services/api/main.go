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
	addr := ":8080"
	topic := os.Getenv("DIRECTIVE_TOPIC")
	if topic == "" {
		topic = "libretto.dev.directive.issued.v1"
	}
	producer := os.Getenv("PRODUCER")
	if producer == "" {
		producer = "api"
	}

	mux := healthMux()
	pub := publisher.NopPublisher{}
	svc := &apiserver.BatonServer{Pub: pub, Topic: topic, Producer: producer}
	mux.Handle(batonv1connect.NewBatonServiceHandler(svc))

	log.Printf("api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
