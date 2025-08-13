package models

import (
	"time"
)

// Session represents a user's login session with refresh token
type Session struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	RefreshToken string    `json:"-" gorm:"not null"` // Hashed refresh token, not returned in JSON
	DeviceID     string    `json:"device_id" gorm:"default:null"`
	UserAgent    string    `json:"user_agent" gorm:"default:null"`
	IPAddress    string    `json:"ip_address" gorm:"default:null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	RevokedAt    time.Time `json:"revoked_at" gorm:"default:null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
