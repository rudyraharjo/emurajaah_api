package models

import "time"

// City struct
type City struct {
	ID         int       `json:"id"`
	ProvinceID int       `json:"province_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RequestIDCity struct
type RequestIDCity struct {
	ID int `json:"id"`
}
