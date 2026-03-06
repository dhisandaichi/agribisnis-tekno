package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `json:"userId"`
	Type      string         `json:"type"` // e.g., "SMS", "IN_APP"
	Message   string         `json:"message"`
	Status    string         `json:"status"` // e.g., "sent", "pending", "failed"
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
