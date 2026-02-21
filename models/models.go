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

type Genre struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `gorm:"type:varchar(100);unique" json:"name"`
}

type Anime struct {
    ID                 uint      `gorm:"primaryKey" json:"id"`
    Title              string    `gorm:"type:varchar(255)" json:"title"`
    Description        string    `gorm:"type:text" json:"description"`
    ReleaseDate        time.Time `json:"release_date"`
    GenreID            *uint     `json:"genre_id"`

    // JPG images
    ImageURL           string    `gorm:"type:varchar(500)" json:"image_url"`
    SmallImageURL      string    `gorm:"type:varchar(500)" json:"small_image_url"`
    LargeImageURL      string    `gorm:"type:varchar(500)" json:"large_image_url"`

    // WebP images
    ImageURLWebP       string    `gorm:"type:varchar(500)" json:"image_url_webp"`
    SmallImageURLWebP  string    `gorm:"type:varchar(500)" json:"small_image_url_webp"`
    LargeImageURLWebP  string    `gorm:"type:varchar(500)" json:"large_image_url_webp"`

    CreatedAt          time.Time `json:"created_at"`
}


type Episode struct {
    ID           uint   `gorm:"primaryKey" json:"id"`
    AnimeID      uint   `gorm:"index" json:"anime_id"`
    EpisodeNum   int    `json:"episode_num"`
    VideoURL     string `gorm:"type:varchar(500)" json:"video_url"`
    Duration     string `gorm:"type:varchar(20)" json:"duration"`
    ThumbnailURL string `gorm:"type:varchar(500)" json:"thumbnail_url"`
}

// Table names
func (User) TableName() string    { return "users" }
func (AuthLog) TableName() string { return "auth_logs" }
func (Setting) TableName() string { return "settings" }
func (Genre) TableName() string   { return "genres" }
func (Anime) TableName() string   { return "anime" }
func (Episode) TableName() string { return "episodes" }
