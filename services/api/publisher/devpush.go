package publisher

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// DevPushPublisher simulates Pub/Sub push by HTTP POSTing to Plot Weaver's /push endpoint.
// It wraps the event bytes into a Pub/Sub push-like envelope with base64-encoded data.
type DevPushPublisher struct {
	URL    string
	Client *http.Client
}

type devPushEnvelope struct {
	Message struct {
		Data       string            `json:"data"`
		Attributes map[string]string `json:"attributes"`
		MessageID  string            `json:"messageId"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func (d DevPushPublisher) client() *http.Client {
	if d.Client != nil {
		return d.Client
	}
	return &http.Client{Timeout: 5 * time.Second}
}

func (d DevPushPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	if d.URL == "" {
		return fmt.Errorf("devpush URL not configured")
	}
	env := devPushEnvelope{}
	enc := base64.StdEncoding.EncodeToString(data)
	env.Message.Data = enc
	env.Message.Attributes = map[string]string{
		"topic": topic,
	}
	env.Message.MessageID = uuid.NewString()
	env.Subscription = "devpush-local"

	b, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal devpush envelope: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.URL, bytesReader(b))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.client().Do(req)
	if err != nil {
		return fmt.Errorf("post devpush: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("devpush unexpected status: %s", resp.Status)
	}
	return nil
}

// bytesReader returns an io.Reader for b without importing bytes in callers.
// Separated to keep imports tidy in publish path.
func bytesReader(b []byte) *bytes.Reader { // small helper to avoid adding bytes import clutter inline
	return bytes.NewReader(b)
}
