package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/barrynorthern/libretto/gen/go/libretto/graph/v1/graphv1connect"
	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/graphwrite"
	gwserver "github.com/barrynorthern/libretto/services/graphwrite/server"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "libretto.db"
	}

	database, err := db.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run migrations
	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize GraphWrite service
	graphWriteService := graphwrite.NewService(database)

	// Initialize HTTP server
	mux := http.NewServeMux()
	svc := gwserver.NewGraphWriteServer(graphWriteService)
	mux.Handle(graphv1connect.NewGraphWriteServiceHandler(svc))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	addr := ":" + port
	log.Printf("graphwrite listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
