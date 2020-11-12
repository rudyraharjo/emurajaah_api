package models

import "time"

// ResponsePoint struct
type ResponsePoint struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Point     int       `json:"point"`
	IsActive  int       `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RequestAddPoint struct
type RequestAddPoint struct {
	Type  string `json:"type"`
	Point int    `json:"point"`
}

type RequestIDPoint struct {
	ID int `json:"id"`
}
