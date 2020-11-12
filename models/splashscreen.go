package models

import "time"

// SplashScreen Struct
type SplashScreen struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Position    int       `json:"position"`
	IsActive    int       `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RequestIDSplashScreen struct
type RequestIDSplashScreen struct {
	ID int `json:"id"`
}
