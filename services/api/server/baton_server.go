package server

import (
	"context"
	"os"

	"connectrpc.com/connect"
	batonv1 "github.com/barrynorthern/libretto/gen/go/libretto/baton/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
	eventsv1 "github.com/barrynorthern/libretto/gen/go/libretto/events/v1"
	contracts_events "github.com/barrynorthern/libretto/packages/contracts/events"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	// Build typed Event with Envelope and oneof payload
	version := os.Getenv("EVENT_VERSION")
	if version == "" {
		version = "1.0.0"
	}
	ev := &eventsv1.Event{
		Envelope: &eventsv1.Envelope{
			EventName:      "DirectiveIssued",
			EventVersion:   version,
			EventId:        uuid.NewString(),
			OccurredAt:     timestamppb.Now(),
			CorrelationId:  uuid.NewString(),
			CausationId:    uuid.NewString(), // non-empty for root events
			IdempotencyKey: uuid.NewString(),
			Producer:       s.Producer,
			TenantId:       "dev",
		},
		Payload: &eventsv1.Event_DirectiveIssued{
			DirectiveIssued: &eventsv1.DirectiveIssued{
				Text:   req.Msg.GetText(),
				Act:    req.Msg.GetAct(),
				Target: req.Msg.GetTarget(),
			},
		},
	}
	b, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(ev)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Optional proto-based envelope validation (default on)
	if os.Getenv("ENVELOPE_VALIDATE") != "0" {
		if err := contracts_events.ValidateEnvelope(ev.GetEnvelope()); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	if err := s.Pub.Publish(ctx, s.Topic, b); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&batonv1.IssueDirectiveResponse{CorrelationId: ev.GetEnvelope().GetCorrelationId()})
	return res, nil
}
