package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/barrynorthern/libretto/gen/go/libretto/baton/v1/batonv1connect"
	"github.com/barrynorthern/libretto/internal/app"
	"github.com/barrynorthern/libretto/internal/db"
	gwpkg "github.com/barrynorthern/libretto/internal/graphwrite"
	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Initialize database
	database, err := db.NewDatabase("libretto.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create a basic project and version for the orchestrator
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	
	// Initialize GraphWrite service
	service := gwpkg.NewService(database)

	mux := http.NewServeMux()

	// Wire orchestrated Baton service
	orchestrator := app.NewOrchestrator(service, versionID)
	mux.Handle(batonv1connect.NewBatonServiceHandler(orchestrator))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("libretto (monolith) listening on %s (project: %s, version: %s)", addr, projectID, versionID)
	log.Fatal(http.ListenAndServe(addr, mux))
}
