package models

import (
	"time"
)

// AuditAction represents the type of authentication action
type AuditAction string

const (
	ActionLogin         AuditAction = "login"
	ActionLogout        AuditAction = "logout"
	ActionRefresh       AuditAction = "refresh"
	ActionPasswordReset AuditAction = "password_reset"
	ActionRegister      AuditAction = "register"
	ActionGoogleAuth    AuditAction = "google_auth"
)

// AuthAudit represents an authentication event for auditing purposes
type AuthAudit struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    uint        `json:"user_id" gorm:"index"`
	Action    AuditAction `json:"action" gorm:"not null"`
	Success   bool        `json:"success" gorm:"default:true"`
	IPAddress string      `json:"ip_address" gorm:"default:null"`
	UserAgent string      `json:"user_agent" gorm:"default:null"`
	DeviceID  string      `json:"device_id" gorm:"default:null"`
	Details   string      `json:"details" gorm:"type:text;default:null"`
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime"`
}
