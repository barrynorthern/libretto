package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/graph/v1/graphv1connect"
)

// InMemoryStore is a minimal stub persistence for Apply.
type InMemoryStore struct{}

func (s *InMemoryStore) Apply(parent string, deltas []*graphv1.Delta) (string, int32, error) {
	// Return a fake new version id and count applied
	return "01JFAKEVERSION", int32(len(deltas)), nil
}

type GraphWriteServer struct {
	graphv1connect.UnimplementedGraphWriteServiceHandler
	Store interface {
		Apply(parent string, deltas []*graphv1.Delta) (string, int32, error)
	}
}

func (s *GraphWriteServer) Apply(ctx context.Context, req *connect.Request[graphv1.ApplyRequest]) (*connect.Response[graphv1.ApplyResponse], error) {
	if len(req.Msg.GetDeltas()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no deltas provided"))
	}
	version, count, err := s.Store.Apply(req.Msg.GetParentVersionId(), req.Msg.GetDeltas())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&graphv1.ApplyResponse{GraphVersionId: version, Applied: count})
	return res, nil
}
