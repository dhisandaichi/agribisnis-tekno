package models

import "time"

type Farmer struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name"`
	Location         string    `json:"location"` // Will hold string representation of GEOGRAPHY
	CropType         string    `json:"crop_type"`
	LandArea         float64   `json:"land_area"`
	FertilizerUsage  float64   `json:"fertilizer_usage"`
	CompetitorProduct string    `json:"competitor_product"`
	CreatedAt        time.Time `json:"created_at"`
}
