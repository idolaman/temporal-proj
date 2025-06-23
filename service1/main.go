package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	// Temporal connection settings
	DefaultTemporalAddress = "localhost:7233"
	DefaultNamespace       = "default"

	// Task queue settings
	TaskQueueName = "url-scanner-task-queue"
)

func main() {
	log.Println("Starting URL Scanner Service (Service #1)")

	// Create Temporal client
	temporalAddress := getEnvOrDefault("TEMPORAL_ADDRESS", DefaultTemporalAddress)
	namespace := getEnvOrDefault("TEMPORAL_NAMESPACE", DefaultNamespace)

	clientOptions := client.Options{
		HostPort:  temporalAddress,
		Namespace: namespace,
	}

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	log.Printf("Connected to Temporal server at %s (namespace: %s)", temporalAddress, namespace)

	w := worker.New(temporalClient, TaskQueueName, worker.Options{})

	scannerWorkflow := &ScannerWorkflow{}
	w.RegisterWorkflow(scannerWorkflow.ScanURLWorkflow)

	scannerActivities := NewScannerActivities()
	w.RegisterActivity(scannerActivities.ScanURL)

	log.Printf("Registered workflow and activity on task queue: %s", TaskQueueName)

	// Start worker
	log.Println("Starting worker...")
	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Worker failed to start: %v", err)
		}
	}()

	log.Println("URL Scanner Service is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Println("Shutting down URL Scanner Service...")
	w.Stop()
	log.Println("URL Scanner Service stopped")
}

// getEnvOrDefault returns the environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 