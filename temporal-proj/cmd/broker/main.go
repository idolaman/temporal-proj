package main

import (
	"context"
	"log"

	acts "temporal-proj/activities"
	temporal "temporal-proj/temporal"
	wflows "temporal-proj/workflows"
)

func main() {
	// Temporal client for scan workflow
	exec, err := temporal.NewClient("ScanURLWorkflow", "url-scanner-task-queue")
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer exec.Close()

	opts := temporal.WorkerOpts{
		TaskQueue:  "url-scanner-task-queue",
		Workflows:  []interface{}{(&wflows.ScannerWorkflow{}).ScanURLWorkflow},
		Activities: []interface{}{acts.NewScannerActivities().ScanURL},
	}

	ctx := context.Background()
	if err := temporal.RunWorker(ctx, exec, opts); err != nil {
		log.Fatalf("Worker exited with error: %v", err)
	}
}
