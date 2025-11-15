package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID    uint           `json:"role_id"`
	Fullname  string         `json:"fullname"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
