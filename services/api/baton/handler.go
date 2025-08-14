package baton

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/barrynorthern/libretto/packages/contracts/events"
	"github.com/google/uuid"
)

type DirectiveRequest struct {
	Text   string `json:"text"`
	Act    string `json:"act,omitempty"`
	Target string `json:"target,omitempty"`
}

type Publisher interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

func Handler(pub Publisher, topic string, producer string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DirectiveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Text == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
			return
		}
		payload := map[string]any{"text": req.Text, "act": req.Act, "target": req.Target}
		ev := events.Envelope[map[string]any]{
			EventName:      "DirectiveIssued",
			EventVersion:   "1.0.0",
			EventID:        uuid.NewString(),
			OccurredAt:     time.Now().UTC(),
			CorrelationID:  uuid.NewString(),
			CausationID:    "",
			IdempotencyKey: uuid.NewString(),
			Producer:       producer,
			TenantID:       "dev",
			Payload:        payload,
		}
		b, _ := json.Marshal(ev)
		if err := pub.Publish(r.Context(), topic, b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "publish failed"})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
	}
}

