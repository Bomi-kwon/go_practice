package model

import (
	"time"
)

// 기본 모델 구조체
type Base struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
