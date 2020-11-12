package models

import "time"

// TermCondition Struct
type TermCondition struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsActived   int       `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
