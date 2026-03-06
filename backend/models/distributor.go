package models

import (
	"time"

	"gorm.io/gorm"
)

type Distributor struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `json:"name"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

type Stock struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	DistributorID uint      `json:"distributorId"`
	ProductName   string    `json:"productName"`
	Quantity      int       `json:"quantity"`
	ReorderLevel  int       `json:"reorderLevel"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Debt struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	DistributorID uint      `json:"distributorId"`
	Amount        float64   `json:"amount"`
	DueDate       time.Time `json:"dueDate"`
	Status        string    `json:"status"` // e.g., "unpaid", "paid"
	CreatedAt     time.Time `json:"createdAt"`
}
