package temporal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/worker"
)

type WorkerOpts struct {
	TaskQueue  string
	Workflows  []interface{}
	Activities []interface{}
}

func RunWorker(ctx context.Context, exec *Client, opts WorkerOpts) error {
	if opts.TaskQueue == "" {
		return fmt.Errorf("task queue must be specified")
	}

	w := worker.New(exec.c, opts.TaskQueue, worker.Options{})

	for _, wf := range opts.Workflows {
		w.RegisterWorkflow(wf)
	}
	for _, act := range opts.Activities {
		w.RegisterActivity(act)
	}

	log.Printf("Registered %d workflows and %d activities on task queue %s", len(opts.Workflows), len(opts.Activities), opts.TaskQueue)

	// Run worker in background so we can capture shutdown signals.
	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalf("Worker runtime error: %v", err)
		}
	}()

	log.Println("Worker is running. Press Ctrl+C to stop.")

	// Listen for cancellation or OS signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down worker...")
	case <-sigCh:
		log.Println("Interrupt received, shutting down worker...")
	}

	w.Stop()
	log.Println("Worker stopped")
	return nil
}
