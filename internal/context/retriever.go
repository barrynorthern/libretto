package contextmgr

// Result represents a retrieval hit with score and optional metadata
// Results are domain-level; UI-facing DTOs are defined in protobufs at the boundary.

type Result struct {
	DocID    string
	Kind     string
	Text     string
	Score    float64 // higher is better for cosine similarity
}

type Retriever interface {
	// Search returns top-k results for the query within a project scope.
	Search(projectID string, query string, k int) ([]Result, error)
}

