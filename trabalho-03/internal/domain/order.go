package domain

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID string         `json:"customer_id"`
	Amount     float64        `json:"amount"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
