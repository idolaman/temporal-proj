package repository

import "time"

// Scan represents a scan request stored in DB.
type Scan struct {
	ID        uint   `gorm:"primaryKey"`
	URL       string `gorm:"not null"`
	Status    string `gorm:"default:pending"`
	LinkCount int    `gorm:"default:0"`
	CreatedAt time.Time
}

// Link represents a discovered link associated with a Scan.
type Link struct {
	ID     uint   `gorm:"primaryKey"`
	ScanID uint   `gorm:"not null"`
	URL    string `gorm:"not null"`
}
