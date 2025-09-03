package main

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SceneDTO mirrors libretto.scene.v1.Scene at the UI boundary.
// In the future we will import the generated proto types directly.
type SceneDTO struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Created string `json:"created"`
}

// App holds minimal in-memory state for a usable scaffold.
type App struct {
	ctx    context.Context
	mu     sync.Mutex
	scenes []*SceneDTO
}

// NewApp creates a new App application struct
func NewApp() *App { return &App{scenes: []*SceneDTO{}} }

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// ListScenes returns current scenes (in-memory for now).
func (a *App) ListScenes() []*SceneDTO {
	a.mu.Lock()
	defer a.mu.Unlock()
	// return a copy to avoid mutation from JS side
	out := make([]*SceneDTO, len(a.scenes))
	copy(out, a.scenes)
	return out
}

// CreateScene adds a new scene and returns it.
func (a *App) CreateScene(title, summary, content string) *SceneDTO {
	a.mu.Lock()
	defer a.mu.Unlock()
	s := &SceneDTO{
		Id:      uuid.NewString(),
		Title:   title,
		Summary: summary,
		Content: content,
		Created: time.Now().Format(time.RFC3339),
	}
	a.scenes = append(a.scenes, s)
	return s
}
