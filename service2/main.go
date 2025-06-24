package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nomaproj/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ScanRequest is the JSON payload sent to the /scan endpoint.
type ScanRequest struct {
	URL string `json:"url" binding:"required"`
}

// Scan represents a row in the "scans" table.
type Scan struct {
	ID        uint   `gorm:"primaryKey"`
	URL       string `gorm:"not null"`
	Status    string `gorm:"default:pending"`
	LinkCount int    `gorm:"default:0"`
	CreatedAt time.Time
}

// Link represents a row in the "links" table.
type Link struct {
	ID     uint   `gorm:"primaryKey"`
	ScanID uint   `gorm:"not null"`
	URL    string `gorm:"not null"`
}

var db *gorm.DB
var temporalClient *TemporalClient

func main() {
	var err error
	dbHost := utils.GetEnvOrDefault("DB_HOST", "localhost")
	dbPort := utils.GetEnvOrDefault("DB_PORT", "5432")
	dbUser := utils.GetEnvOrDefault("DB_USER", "postgres")
	dbPassword := utils.GetEnvOrDefault("DB_PASSWORD", "password")
	dbName := utils.GetEnvOrDefault("DB_NAME", "crawler")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Migrate database
	db.AutoMigrate(&Scan{}, &Link{})

	// Connect to Temporal
	temporalClient, err = NewTemporalClient()
	if err != nil {
		log.Fatal("Temporal connection failed:", err)
	}
	defer temporalClient.Close()

	// Setup API
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/scan", handleScan)

	port := utils.GetEnvOrDefault("PORT", "8080")
	log.Printf("Service #2 running on :%s", port)
	log.Fatal(router.Run(":" + port))
}

func handleScan(c *gin.Context) {
	var request ScanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Create scan record
	scan := Scan{URL: request.URL, Status: "pending"}
	if err := db.Create(&scan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 2. Start Service #1 workflow
	if err := temporalClient.StartScan(c.Request.Context(), request.URL, scan.ID); err != nil {
		db.Model(&scan).Update("status", "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start scan"})
		return
	}

	// 3. Start background process to wait for results
	go waitForResults(scan.ID)

	c.JSON(http.StatusAccepted, gin.H{
		"id":      scan.ID,
		"url":     scan.URL,
		"status":  "pending",
		"message": "Scan started",
	})
}

func waitForResults(scanID uint) {
	ctx := context.Background()

	// Wait for workflow result
	result, err := temporalClient.GetScanResult(ctx, scanID)
	if err != nil {
		log.Printf("Failed to get result for scan %d: %v", scanID, err)
		db.Model(&Scan{}).Where("id = ?", scanID).Update("status", "failed")
		return
	}

	// Save results to database
	if result.Success {
		tx := db.Begin()

		// Save links
		for _, linkURL := range result.Links {
			link := Link{ScanID: scanID, URL: linkURL}
			if err := tx.Create(&link).Error; err != nil {
				tx.Rollback()
				log.Printf("Failed to save links for scan %d: %v", scanID, err)
				db.Model(&Scan{}).Where("id = ?", scanID).Update("status", "failed")
				return
			}
		}

		// Update scan status
		tx.Model(&Scan{}).Where("id = ?", scanID).Updates(map[string]interface{}{
			"status":     "completed",
			"link_count": len(result.Links),
		})

		tx.Commit()
		log.Printf("Saved %d links for scan %d", len(result.Links), scanID)
	} else {
		db.Model(&Scan{}).Where("id = ?", scanID).Update("status", "failed")
		log.Printf("Scan %d failed: %s", scanID, result.Error)
	}
}
