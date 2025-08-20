package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPushHandlerRejectsInvalidEvent(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/push", bytes.NewReader([]byte(`{"message":{"data":"not-base64"}}`)))
	pushHandler(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPushHandlerAcceptsValidEvent(t *testing.T) {
	// Build a minimal valid Event JSON and base64 it
	ev := &eventsv1.Event{Envelope: &eventsv1.Envelope{EventName: "DirectiveIssued", EventVersion: "1.0.0", EventId: "id", CorrelationId: "corr", CausationId: "cause", IdempotencyKey: "idem", Producer: "plotweaver", TenantId: "dev", OccurredAt: timestamppb.Now()}}
	b, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(ev)
	enc := base64.StdEncoding.EncodeToString(b)
	body := map[string]any{"message": map[string]any{"data": enc, "attributes": map[string]string{}, "messageId": "1"}, "subscription": "devpush"}
	raw, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/push", bytes.NewReader(raw))
	pushHandler(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

