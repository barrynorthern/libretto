package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
	"github.com/barrynorthern/libretto/internal/graphwrite"
)

func TestApplyRejectsEmptyDeltas(t *testing.T) {
	s := NewGraphWriteServer(&mockGraphWriteService{version: "01JF00", count: 0})
	req := connect.NewRequest(&graphv1.ApplyRequest{ParentVersionId: "01JROOT", Deltas: []*graphv1.Delta{}})
	_, err := s.Apply(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error for empty deltas")
	}
	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Fatalf("expected invalid argument, got %v", connect.CodeOf(err))
	}
}

type mockGraphWriteService struct {
	version string
	count   int32
	err     error
}

func (m *mockGraphWriteService) Apply(ctx context.Context, req *graphwrite.ApplyRequest) (*graphwrite.ApplyResponse, error) {
	if len(req.Deltas) == 0 {
		return nil, m.err
	}
	return &graphwrite.ApplyResponse{
		GraphVersionID: m.version,
		Applied:        m.count,
	}, m.err
}

func (m *mockGraphWriteService) GetVersion(ctx context.Context, versionID string) (*graphwrite.GraphVersion, error) {
	return nil, m.err
}

func (m *mockGraphWriteService) ListEntities(ctx context.Context, versionID string, filter graphwrite.EntityFilter) ([]*graphwrite.Entity, error) {
	return nil, m.err
}

func (m *mockGraphWriteService) GetNeighbors(ctx context.Context, entityID string, relationshipType string) ([]*graphwrite.Entity, error) {
	return nil, m.err
}

func (m *mockGraphWriteService) GetNeighborsInVersion(ctx context.Context, versionID string, logicalEntityID string, relationshipType string) ([]*graphwrite.Entity, error) {
	return nil, m.err
}

func (m *mockGraphWriteService) ImportEntity(ctx context.Context, targetVersionID, sourceProjectID, entityLogicalID string) (*graphwrite.Entity, error) {
	return nil, m.err
}

func (m *mockGraphWriteService) GetEntityHistory(ctx context.Context, entityLogicalID string) ([]*graphwrite.EntityVersion, error) {
	return nil, m.err
}

func (m *mockGraphWriteService) ListSharedEntities(ctx context.Context) ([]*graphwrite.SharedEntity, error) {
	return nil, m.err
}

func TestApplySuccess(t *testing.T) {
	s := NewGraphWriteServer(&mockGraphWriteService{version: "01JF00", count: 2})
	req := connect.NewRequest(&graphv1.ApplyRequest{ParentVersionId: "01JROOT", Deltas: []*graphv1.Delta{{Op: "create"}, {Op: "create"}}})
	res, err := s.Apply(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got, want := res.Msg.GetGraphVersionId(), "01JF00"; got != want {
		t.Fatalf("version got %q want %q", got, want)
	}
	if got, want := res.Msg.GetApplied(), int32(2); got != want {
		t.Fatalf("applied got %d want %d", got, want)
	}
}
