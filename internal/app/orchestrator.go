package app

import (
	"context"

	"connectrpc.com/connect"
	batonv1 "github.com/barrynorthern/libretto/gen/go/libretto/baton/v1"
	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
	"github.com/barrynorthern/libretto/internal/agents/narrative"
	"github.com/barrynorthern/libretto/internal/agents/plotweaver"
	gwpkg "github.com/barrynorthern/libretto/internal/graphwrite"
)

// Orchestrator implements BatonService and synchronously calls agent modules.
type Orchestrator struct {
	plot      plotweaver.Module
	narr      narrative.Module
	gw        gwpkg.GraphWriteService
	versionID string
	producer  string
}

func NewOrchestrator(service gwpkg.GraphWriteService, versionID string) *Orchestrator {
	return &Orchestrator{
		plot:      plotweaver.New(),
		narr:      narrative.New(),
		gw:        service,
		versionID: versionID,
		producer:  "monolith",
	}
}

var _ batonv1connect.BatonServiceHandler = (*Orchestrator)(nil)

func (o *Orchestrator) IssueDirective(ctx context.Context, req *connect.Request[batonv1.IssueDirectiveRequest]) (*connect.Response[batonv1.IssueDirectiveResponse], error) {
	// Synchronously process directive
	proposal := o.plot.ProcessDirective(ctx, req.Msg.GetText(), req.Msg.GetAct(), req.Msg.GetTarget(), o.producer)
	// Apply to store using the current GraphWrite service
	_ = o.narr.ApplySceneProposal(ctx, o.gw, o.versionID, proposal)
	return connect.NewResponse(&batonv1.IssueDirectiveResponse{CorrelationId: proposal.CorrelationId}), nil
}
