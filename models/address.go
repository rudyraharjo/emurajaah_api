package models

import "time"

// Province struct
type Province struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	InternationID string    `json:"internation_id"`
	CountryID     int       `json:"country_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// PropKey       string    `json:"prop_key"`
	// PropKeyOrg    string    `json:"prop_key_org"`
	// PropName      string    `json:"prop_name"`
	// CreationDate  time.Time `json:"creation_date"`
	// InternationId string    `json:"internation_id"`
	// Ibukota       string    `json:"ibukota"`
	// Korwil        string    `json:"korwil"`
}

// RequestIDProvince struct
type RequestIDProvince struct {
	ID int `json:"id"`
}
