package models

import "time"

type Farmer struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	CropType          string    `json:"crop_type"`
	LandArea          float64   `json:"land_area"`
	FertilizerUsage   float64   `json:"fertilizer_usage"`
	CompetitorProduct string    `json:"competitor_product"`
	CreatedAt         time.Time `json:"created_at"`
}
