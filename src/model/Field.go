package models

import (
	"time"

	"gorm.io/gorm"
)

type Field struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Price     int            `json:"price"`
	Location  string         `json:"location"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
