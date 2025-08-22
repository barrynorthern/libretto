package graphwrite

import "context"

type Store interface {
	CreateScene(ctx context.Context, id, title, summary string) error
}

func NewInMemory() Store { return &mem{} }

type mem struct{}

func (m *mem) CreateScene(ctx context.Context, id, title, summary string) error {
	_ = ctx; _ = id; _ = title; _ = summary
	return nil
}

