package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Role         string    `json:"role"` // 'sales', 'supervisor', 'farmer', 'distributor'
	Email        string    `json:"email" gorm:"unique"`
	PasswordHash string    `json:"-"` // Don't expose this fields in JSON
	CreatedAt    time.Time `json:"created_at"`
}
