package models

import "time"

// Country struct
type Country struct {
	ID          int       `json:"id"`
	CountryName string    `json:"country_name"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RequestIDCountry struct
type RequestIDCountry struct {
	ID int `json:"id"`
}
