package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type pushEnvelope struct {
	Message struct {
		Data       string            `json:"data"`
		Attributes map[string]string `json:"attributes"`
		MessageID  string            `json:"messageId"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// pushHandler accepts Pub/Sub push messages. For now, it logs the decoded event envelope
// and responds 200. It coexists with the existing root handler used in local stub flows.
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
	log.Printf("plotweaver: received push messageId=%s attrs=%v payload=%s", env.Message.MessageID, env.Message.Attributes, string(dec))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

type PubSubMessage struct {
	Message struct {
		Data       []byte            `json:"data"`
		ID         string            `json:"messageId"`
		Attributes map[string]string `json:"attributes"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func publishSceneProposal(ctx context.Context, topic string, publish func(context.Context, string, []byte) error) error {
	ev := map[string]any{
		"eventName":      "SceneProposalReady",
		"eventVersion":   "1.0.0",
		"eventId":        uuid.NewString(),
		"occurredAt":     time.Now().UTC().Format(time.RFC3339Nano),
		"correlationId":  uuid.NewString(),
		"causationId":    "",
		"idempotencyKey": uuid.NewString(),
		"producer":       "plotweaver",
		"tenantId":       "dev",
		"payload": map[string]any{
			"scene_id": uuid.NewString(),
			"title":    "A turning point",
			"summary":  "A betrayal changes the course of events.",
		},
	}
	b, _ := json.Marshal(ev)
	return publish(ctx, topic, b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// In MVP, we ignore the contents and always emit a stub proposal
	ctx := r.Context()
	_ = publishSceneProposal(ctx, "libretto.dev.scene.proposal.ready.v1", func(ctx context.Context, topic string, data []byte) error {
		log.Printf("(stub) publish to %s: %s", topic, string(data))
		return nil
	})
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
