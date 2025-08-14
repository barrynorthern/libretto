package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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

