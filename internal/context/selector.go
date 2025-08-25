package contextmgr

// ModelSpec describes a selected model configuration.

type ModelSpec struct {
	Provider string // "ollama" or provider key
	Model    string // e.g., "llama3.1:8b-instruct"
	MaxTokens int
	BudgetUSD float64 // per-call budget ceiling
}

type ModelSelector interface {
	// Choose selects an appropriate model for a task given complexity and budget.
	Choose(task string, complexity string, budgetUSD float64) (ModelSpec, error)
}

