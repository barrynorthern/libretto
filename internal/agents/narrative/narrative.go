package narrative

import (
	"context"
	"log"

	"github.com/barrynorthern/libretto/internal/agents/plotweaver"
	"github.com/barrynorthern/libretto/internal/graphwrite"
)

type Module interface {
	ApplySceneProposal(ctx context.Context, service graphwrite.GraphWriteService, versionID string, p plotweaver.SceneProposal) error
}

func New() Module { return &impl{} }

type impl struct{}

func (i *impl) ApplySceneProposal(ctx context.Context, service graphwrite.GraphWriteService, versionID string, p plotweaver.SceneProposal) error {
	// Map proposal to a graph delta using the current GraphWrite service
	req := &graphwrite.ApplyRequest{
		ParentVersionID: versionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   p.SceneID,
				Fields: map[string]any{
					"name":        p.Title,
					"description": p.Summary,
				},
			},
		},
	}
	
	_, err := service.Apply(ctx, req)
	if err != nil {
		return err
	}
	
	log.Printf("narrative: applied Scene %s title=%q", p.SceneID, p.Title)
	return nil
}

