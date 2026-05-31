package authmodel

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"size:64;not null;uniqueIndex"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;size:255;not null"`
	Name         string         `json:"name" gorm:"size:128;not null"`
	FlagActive   bool           `json:"flag_active" gorm:"not null;default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
