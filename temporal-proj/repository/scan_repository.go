package repository

import (
	"context"
	"fmt"

	"temporal-proj/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ScanRepository interface {
	CreateScan(ctx context.Context, s *Scan) error
	UpdateStatus(ctx context.Context, scanID uint, status string, linkCount int) error
	AddLinks(ctx context.Context, scanID uint, links []Link) error
}

func PostgresScanRepository() (ScanRepository, error) {
	// connection details are fetched here so that callers remain agnostic.
	dbHost := utils.GetEnvOrDefault("DB_HOST", "localhost")
	dbPort := utils.GetEnvOrDefault("DB_PORT", "5432")
	dbUser := utils.GetEnvOrDefault("DB_USER", "postgres")
	dbPassword := utils.GetEnvOrDefault("DB_PASSWORD", "password")
	dbName := utils.GetEnvOrDefault("DB_NAME", "crawler")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Scan{}, &Link{}); err != nil {
		return nil, err
	}

	return &postgresScanRepo{db: db}, nil
}


type postgresScanRepo struct {
	db *gorm.DB
}

func (r *postgresScanRepo) CreateScan(ctx context.Context, s *Scan) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *postgresScanRepo) UpdateStatus(ctx context.Context, scanID uint, status string, linkCount int) error {
	updates := map[string]interface{}{"status": status}
	if linkCount >= 0 {
		updates["link_count"] = linkCount
	}
	return r.db.WithContext(ctx).Model(&Scan{}).Where("id = ?", scanID).Updates(updates).Error
}

func (r *postgresScanRepo) AddLinks(ctx context.Context, scanID uint, links []Link) error {
	if len(links) == 0 {
		return nil
	}
	// Use transaction for consistency
	tx := r.db.WithContext(ctx).Begin()
	for _, l := range links {
		if err := tx.Create(&l).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
