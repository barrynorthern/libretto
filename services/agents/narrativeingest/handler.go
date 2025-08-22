package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/graph/v1/graphv1connect"
	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type pushEnvelope struct {
	Message struct {
		Data       string            `json:"data"`
		Attributes map[string]string `json:"attributes"`
		MessageID  string            `json:"messageId"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

var (
	graphClient graphv1connect.GraphWriteServiceClient
)

func init() {
	url := os.Getenv("GRAPHWRITE_URL")
	if url == "" {
		url = "http://localhost:8082"
	}
	graphClient = graphv1connect.NewGraphWriteServiceClient(http.DefaultClient, url)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}
	var env pushEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		http.Error(w, "invalid push envelope", http.StatusBadRequest)
		return
	}
	dec, err := base64.StdEncoding.DecodeString(env.Message.Data)
	if err != nil {
		http.Error(w, "invalid base64 data", http.StatusBadRequest)
		return
	}
	var ev eventsv1.Event
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(dec, &ev); err != nil {
		log.Printf("narrative-ingest: event decode error: %v", err)
		http.Error(w, "invalid event", http.StatusBadRequest)
		return
	}
	switch p := ev.Payload.(type) {
	case *eventsv1.Event_SceneProposalReady:
		sp := p.SceneProposalReady
		// Map to GraphWrite delta
		req := connect.NewRequest(&graphv1.ApplyRequest{
			ParentVersionId: "01JROOT",
			Deltas: []*graphv1.Delta{{
				Op:         "create",
				EntityType: "Scene",
				EntityId:   sp.GetSceneId(),
				Fields:     map[string]string{"title": sp.GetTitle(), "summary": sp.GetSummary()},
			}},
		})
		if _, err := graphClient.Apply(r.Context(), req); err != nil {
			log.Printf("narrative-ingest: graph apply error: %v", err)
			http.Error(w, "graph apply error", http.StatusBadRequest)
			return
		}
		log.Printf("narrative-ingest: applied Scene %s title=%q", sp.GetSceneId(), sp.GetTitle())
	default:
		// ignore other events
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

