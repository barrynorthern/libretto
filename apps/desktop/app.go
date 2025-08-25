package main

import (
	"context"
)

// SceneDTO mirrors libretto.scene.v1.Scene at the UI boundary.
// In the future we will import the generated proto types directly.
type SceneDTO struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ListScenes returns an array of scene DTOs.
// For now, return an empty list; will wire to repository later.
func (a *App) ListScenes() []*SceneDTO {
	return []*SceneDTO{}
}
