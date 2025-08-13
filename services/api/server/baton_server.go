package server

import (
	"context"
	"encoding/json"
	"time"

	"connectrpc.com/connect"
	batonv1 "github.com/barrynorthern/libretto/gen/go/libretto/baton/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	"github.com/google/uuid"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

type BatonServer struct {
	batonv1connect.UnimplementedBatonServiceHandler
	Pub      Publisher
	Topic    string
	Producer string
}

func (s *BatonServer) IssueDirective(ctx context.Context, req *connect.Request[batonv1.IssueDirectiveRequest]) (*connect.Response[batonv1.IssueDirectiveResponse], error) {
	env := map[string]any{
		"eventName":      "DirectiveIssued",
		"eventVersion":   "1.0.0",
		"eventId":        uuid.NewString(),
		"occurredAt":     time.Now().UTC().Format(time.RFC3339Nano),
		"correlationId":  uuid.NewString(),
		"causationId":    "",
		"idempotencyKey": uuid.NewString(),
		"producer":       s.Producer,
		"tenantId":       "dev",
		"payload": map[string]any{
			"text":   req.Msg.GetText(),
			"act":    req.Msg.GetAct(),
			"target": req.Msg.GetTarget(),
		},
	}
	b, _ := json.Marshal(env)
	if err := s.Pub.Publish(ctx, s.Topic, b); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&batonv1.IssueDirectiveResponse{CorrelationId: env["correlationId"].(string)})
	return res, nil
}

