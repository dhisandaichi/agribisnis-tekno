package models

import (
	"time"

	"gorm.io/gorm"
)

type Content struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	VideoURL    string         `json:"videoUrl"`
	Thumbnail   string         `json:"thumbnail"`
	Category    string         `json:"category"` // e.g., "education", "market_price"
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
