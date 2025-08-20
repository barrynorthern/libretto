package server

import (
	"encoding/json"
	"testing"

	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestEventMarshallingOneofDirectiveIssued(t *testing.T) {
	ev := &eventsv1.Event{
		Envelope: &eventsv1.Envelope{EventName: "DirectiveIssued", EventVersion: "1.0.0"},
		Payload:  &eventsv1.Event_DirectiveIssued{DirectiveIssued: &eventsv1.DirectiveIssued{Text: "x"}},
	}
	b, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if _, ok := m["directiveIssued"]; !ok {
		t.Fatalf("expected directiveIssued oneof field present in JSON")
	}
}

