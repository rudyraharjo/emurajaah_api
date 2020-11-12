package models

import "time"

// Banner Struct
type Banner struct {
	Id             int       `json:"id"`
	Title          string    `json:"title"`
	Subtitle       string    `json:"subtitle"`
	ImageUrl       string    `json:"image_url"`
	BannerPosition int       `json:"banner_position"`
	IsActive       int       `json:"is_active"`
	CreatedDate    time.Time `json:"created_date"`
}

// RequestAddBanner Struct
type RequestAddBanner struct {
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	ImageURL       string `json:"image_url"`
	BannerPosition int    `json:"banner_position"`
}

// RequestDeleteBanner Struct
type RequestIDBanner struct {
	ID int `json:"id"`
}
