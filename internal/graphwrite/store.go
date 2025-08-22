package graphwrite

import (
	"context"
	"net/http"
)

type Store interface {
	CreateScene(ctx context.Context, id, title, summary string) error
	ListScenes(r *http.Request) any
}

type MemStore = mem

func NewInMemory() Store { return &mem{} }

type mem struct{}

func (m *mem) CreateScene(ctx context.Context, id, title, summary string) error {
	_ = ctx
	_ = id
	_ = title
	_ = summary
	return nil
}

func (m *mem) ListScenes(r *http.Request) any { _ = r; return []any{} }
