package contextmgr

// Domain types used by the Context Manager

type ContextBundle struct {
	// TODO: add fields like RelevantScenes, Characters, TokenBudget, etc.
}

type Directive struct {
	Text   string
	Act    string
	Target string
}

type Builder interface {
	// Build assembles a ContextBundle for a given project and directive.
	Build(projectID string, directive Directive) (ContextBundle, error)
}

// NoOpBuilder returns an empty bundle; used to wire seams during MVP bring-up
// without introducing external dependencies.
type NoOpBuilder struct{}

func NewNoOpBuilder() *NoOpBuilder { return &NoOpBuilder{} }

func (b *NoOpBuilder) Build(projectID string, directive Directive) (ContextBundle, error) {
	return ContextBundle{}, nil
}

