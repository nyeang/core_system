package models

import (
    "time"
    "gorm.io/gorm"
)

// User model
type User struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Username     string         `gorm:"type:varchar(255);not null" json:"username"`
    Email        string         `gorm:"type:varchar(255);unique;not null" json:"email"`
    PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
    Role         string         `gorm:"type:varchar(50);default:'user'" json:"role"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// AuthLog model - tracks all authentication activities
type AuthLog struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    UserID    uint           `gorm:"index" json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Action    string         `gorm:"type:varchar(50);not null" json:"action"`
    IPAddress string         `gorm:"type:varchar(45)" json:"ip_address"`
    UserAgent string         `gorm:"type:varchar(500)" json:"user_agent"`
    Subsystem string         `gorm:"type:varchar(50)" json:"subsystem"`
    Status    string         `gorm:"type:varchar(20)" json:"status"`
    Details   string         `gorm:"type:text" json:"details"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Setting model - system settings
type Setting struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Key       string         `gorm:"type:varchar(100);unique;not null" json:"key"`
    Value     string         `gorm:"type:text" json:"value"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Table names
func (User) TableName() string    { return "users" }
func (AuthLog) TableName() string { return "auth_logs" }
func (Setting) TableName() string { return "settings" }
