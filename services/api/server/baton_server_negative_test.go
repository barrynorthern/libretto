package server

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"connectrpc.com/connect"
	batonv1 "github.com/barrynorthern/libretto/gen/go/libretto/baton/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
)

type noopPublisher struct{}

func (noopPublisher) Publish(ctx context.Context, topic string, data []byte) error { return nil }

func TestIssueDirectiveRejectsInvalidEnvelopeVersion(t *testing.T) {
	// Force invalid semver to trigger validation error
	t.Setenv("EVENT_VERSION", "badversion")
	svc := &BatonServer{Pub: noopPublisher{}, Topic: "t", Producer: "api"}
	h := batonv1connect.NewBatonServiceHandler(svc)
	r := connect.NewRequest(&batonv1.IssueDirectiveRequest{Text: "x"})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/libretto.baton.v1.BatonService/IssueDirective", newBytesReader(mustJSON(r))))
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// Helpers for building request bodies without adding extra imports in test
func mustJSON(req *connect.Request[batonv1.IssueDirectiveRequest]) []byte {
	b, err := req.MarshalJSON()
	if err != nil { panic(err) }
	return b
}

type bytesReaderT struct{ b []byte }

func newBytesReader(b []byte) *bytesReaderT { return &bytesReaderT{b: b} }
func (r *bytesReaderT) Read(p []byte) (int, error) { n := copy(p, r.b); r.b = r.b[n:]; if n == 0 { return 0, os.ErrClosed } ; return n, nil }

