package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/api"
	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/k8s"
)

func main() {
	log.Println("Starting cluster-viz backend...")

	// Create Kubernetes clients
	clientset, dynamicClient, err := k8s.NewClients()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clients: %v", err)
	}
	log.Println("Kubernetes clients created")

	// Create graph builder
	graphBuilder := k8s.NewGraphBuilder()

	// Create resource watcher
	watcher := k8s.NewResourceWatcher(clientset, dynamicClient, graphBuilder)

	// Create API handler
	handler := api.NewHandler(graphBuilder, watcher)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watcher
	watcher.Start(ctx)

	// Start broadcast loop
	handler.Start(ctx)

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down gracefully...")

	// Stop watcher
	watcher.Stop()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Shutdown complete")
}
