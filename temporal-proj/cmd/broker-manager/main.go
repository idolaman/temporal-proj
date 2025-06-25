package main

import (
	"log"

	"temporal-proj/pkg/utils"

	api "temporal-proj/api"
	repo "temporal-proj/repository"
	svc "temporal-proj/service"
	temporal "temporal-proj/workflowmgmt/temporal"
)

func main() {
	postgresScanRepo, err := repo.PostgresScanRepository()
	if err != nil {
		log.Fatal("Repository initialization failed:", err)
	}

	temporalClient, err := temporal.NewClient("ScanURLWorkflow", "url-scanner-task-queue")
	if err != nil {
		log.Fatal("Temporal connection failed:", err)
	}
	defer temporalClient.Close()

	coordinator := svc.NewCoordinator(postgresScanRepo, temporalClient)

	// Create and run API server
	port := utils.GetEnvOrDefault("PORT", "8080")
	server := api.NewServer(coordinator)
	log.Printf("Broker Manager running on :%s", port)
	log.Fatal(server.Run(":" + port))
}
