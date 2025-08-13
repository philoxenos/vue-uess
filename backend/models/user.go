package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Role represents a user role in the system
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleInspector Role = "inspector"
)

// Roles is a slice of Role that can be stored in the database
type Roles []Role

// Scan implements the sql.Scanner interface for Roles
func (r *Roles) Scan(value interface{}) error {
	if value == nil {
		*r = Roles{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Roles")
	}

	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface for Roles
func (r Roles) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}

	return json.Marshal(r)
}

// User represents a user in the system
type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Email            string    `json:"email" gorm:"unique;not null"`
	Password         string    `json:"-" gorm:"default:null"` // Password not returned in JSON, can be null for Google-only accounts
	GoogleID         string    `json:"google_id" gorm:"unique;index"`
	GoogleSub        string    `json:"google_sub" gorm:"unique;index"` // Subject identifier from Google
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	HasLocalPassword bool      `json:"has_local_password" gorm:"default:false"`
	Roles            Roles     `json:"roles" gorm:"type:json;default:'[\"user\"]'"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	IsAdmin          bool      `json:"is_admin" gorm:"default:false"`
	LastLogin        time.Time `json:"last_login" gorm:"default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
