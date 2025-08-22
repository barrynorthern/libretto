package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type fakeGraphClient struct{}

func (fakeGraphClient) Apply(ctx context.Context, req *connect.Request[graphv1.ApplyRequest]) (*connect.Response[graphv1.ApplyResponse], error) {
	return connect.NewResponse(&graphv1.ApplyResponse{GraphVersionId: "01JFAKE", Applied: 1}), nil
}

func TestPushAcceptsSceneProposalReady(t *testing.T) {
	// Build a minimal Event JSON and base64 it with SceneProposalReady
	sp := &eventsv1.Event{Envelope: &eventsv1.Envelope{EventName: "SceneProposalReady", EventVersion: "1.0.0"}, Payload: &eventsv1.Event_SceneProposalReady{SceneProposalReady: &eventsv1.SceneProposalReady{SceneId: "sc-1", Title: "A", Summary: "B"}}}
	b, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(sp)
	enc := base64.StdEncoding.EncodeToString(b)
	body := map[string]any{"message": map[string]any{"data": enc, "attributes": map[string]string{}, "messageId": "1"}, "subscription": "devpush"}
	raw, _ := json.Marshal(body)
	// Inject fake graph client so we don't make network calls
	graphClient = fakeGraphClient{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/push", bytes.NewReader(raw))
	pushHandler(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
