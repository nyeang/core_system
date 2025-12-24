package models

import (
	"time"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(255);not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" hides password in API responses
	Role         string    `gorm:"type:varchar(50);default:'user'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Genre model
type Genre struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	Animes    []Anime   `gorm:"foreignKey:GenreID" json:"animes,omitempty"` // One-to-many relationship
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Anime model
type Anime struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	GenreID     uint      `gorm:"index" json:"genre_id"`
	Genre       Genre     `gorm:"foreignKey:GenreID" json:"genre,omitempty"` // Belongs to Genre
	Studio      string    `gorm:"type:varchar(255)" json:"studio"`
	Status      string    `gorm:"type:varchar(50);default:'ongoing'" json:"status"` // ongoing, completed, upcoming
	CoverImage  string    `gorm:"type:varchar(500)" json:"cover_image"`
	Episodes    []Episode `gorm:"foreignKey:AnimeID" json:"episodes,omitempty"` // One-to-many relationship
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Episode model
type Episode struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AnimeID      uint      `gorm:"index;not null" json:"anime_id"`
	Anime        Anime     `gorm:"foreignKey:AnimeID" json:"anime,omitempty"` // Belongs to Anime
	EpisodeNum   int       `gorm:"not null" json:"episode_num"`
	VideoURL     string    `gorm:"type:varchar(500)" json:"video_url"`
	Duration     string    `gorm:"type:varchar(20)" json:"duration"` // e.g., "24:30" or use int for seconds
	ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides (optional - GORM will pluralize automatically)
func (User) TableName() string {
	return "users"
}

func (Genre) TableName() string {
	return "genres"
}

func (Anime) TableName() string {
	return "anime" // Keep singular as per your DBML
}

func (Episode) TableName() string {
	return "episode" // Keep singular as per your DBML
}