package plotweaver

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Module interface {
	ProcessDirective(ctx context.Context, text, act, target, producer string) SceneProposal
}

type SceneProposal struct {
	SceneID       string
	Title         string
	Summary       string
	CorrelationId string
	OccurredAt    time.Time
}

func New() Module { return &impl{} }

type impl struct{}

func (i *impl) ProcessDirective(_ context.Context, text, act, target, producer string) SceneProposal {
	_ = text; _ = act; _ = target; _ = producer
	return SceneProposal{
		SceneID:       uuid.NewString(),
		Title:         "A turning point",
		Summary:       "A betrayal changes the course of events.",
		CorrelationId: uuid.NewString(),
		OccurredAt:    time.Now().UTC(),
	}
}

