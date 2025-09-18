package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/graph/v1/graphv1connect"
	"github.com/barrynorthern/libretto/internal/graphwrite"
)

type GraphWriteServer struct {
	graphv1connect.UnimplementedGraphWriteServiceHandler
	service graphwrite.GraphWriteService
}

// NewGraphWriteServer creates a new GraphWriteServer instance
func NewGraphWriteServer(service graphwrite.GraphWriteService) *GraphWriteServer {
	return &GraphWriteServer{
		service: service,
	}
}

func (s *GraphWriteServer) Apply(ctx context.Context, req *connect.Request[graphv1.ApplyRequest]) (*connect.Response[graphv1.ApplyResponse], error) {
	if len(req.Msg.GetDeltas()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no deltas provided"))
	}

	// Convert protobuf deltas to internal format
	deltas := make([]*graphwrite.Delta, len(req.Msg.GetDeltas()))
	for i, pbDelta := range req.Msg.GetDeltas() {
		// Convert fields map from map[string]string to map[string]any
		fields := make(map[string]any)
		for k, v := range pbDelta.GetFields() {
			fields[k] = v
		}

		deltas[i] = &graphwrite.Delta{
			Operation:  pbDelta.GetOp(),
			EntityType: pbDelta.GetEntityType(),
			EntityID:   pbDelta.GetEntityId(),
			Fields:     fields,
		}
	}

	// Apply deltas using the service
	response, err := s.service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: req.Msg.GetParentVersionId(),
		Deltas:          deltas,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&graphv1.ApplyResponse{
		GraphVersionId: response.GraphVersionID,
		Applied:        response.Applied,
	})
	return res, nil
}
