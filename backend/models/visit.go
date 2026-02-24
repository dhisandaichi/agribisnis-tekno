package models

import "time"

type Visit struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SalesID   uint      `json:"sales_id"`
	FarmerID  uint      `json:"farmer_id"`
	PhotoPath string    `json:"photo_path"`
	VisitTime time.Time `json:"visit_time"`
	Notes     string    `json:"notes"`
}
