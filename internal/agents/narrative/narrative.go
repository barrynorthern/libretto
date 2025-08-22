package narrative

import (
	"context"
	"log"

	"github.com/barrynorthern/libretto/internal/agents/plotweaver"
	"github.com/barrynorthern/libretto/internal/graphwrite"
)

type Module interface {
	ApplySceneProposal(ctx context.Context, store graphwrite.Store, p plotweaver.SceneProposal) error
}

func New() Module { return &impl{} }

type impl struct{}

func (i *impl) ApplySceneProposal(ctx context.Context, store graphwrite.Store, p plotweaver.SceneProposal) error {
	// Map proposal to a graph delta; in-memory store just logs and returns success
	if err := store.CreateScene(ctx, p.SceneID, p.Title, p.Summary); err != nil {
		return err
	}
	log.Printf("narrative: applied Scene %s title=%q", p.SceneID, p.Title)
	return nil
}

