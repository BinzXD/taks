package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID        uint           `json:"id"`
	FieldID   uint           `json:"field_id"`
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	CreatedBy uint           `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
