package main

import (
	"context"
	"log"

	acts "temporal-proj/activities"
	temporal "temporal-proj/workflowmgmt/temporal"
	wflows "temporal-proj/workflows"
)

func main() {
	client, err := temporal.NewClient("ScanURLWorkflow", "url-scanner-task-queue")
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer client.Close()

	temporalClient, ok := client.(*temporal.Client)
	if !ok {
		log.Fatalf("Failed to get temporal client")
	}

	opts := temporal.WorkerOpts{
		TaskQueue:  "url-scanner-task-queue",
		Workflows:  []interface{}{(&wflows.ScannerWorkflow{}).ScanURLWorkflow},
		Activities: []interface{}{acts.NewScannerActivities().ScanURL},
	}

	ctx := context.Background()
	if err := temporal.RunWorker(ctx, temporalClient, opts); err != nil {
		log.Fatalf("Worker exited with error: %v", err)
	}
}
