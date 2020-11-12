package models

import "time"

// Quote struct
type Quote struct {
	Id        int       `json:"id"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	IsActive  int       `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RequestIDQuote struct
type RequestIDQuote struct {
	ID int `json:"id"`
}
