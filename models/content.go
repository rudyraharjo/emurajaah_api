package models

import "time"

// ResponseHomePageContent Struct
type ResponseHomePageContent struct {
	Banner              []Banner                  `json:"banner"`
	Groups              []ResponseGroupUserJoined `json:"groups"`
	Quotes              []Quote                   `json:"quotes"`
	ResponseTotalKhatam []ResponseTotalKhatam     `json:"global_group_status"`
}

// type ResponseHomePageContent struct {
// 	Banner               []Banner               `json:"banner"`
// 	Groups               []ResponseGroupList    `json:"groups"`
// 	Quotes               []Quote                `json:"quotes"`
// 	StatisticOfReadQuran []StatisticOfReadQuran `json:"global_group_status"`
// 	TotalKhatam          *ResponseTotalKhatam   `json:"global_allgroup_total_khatam"`
// }

// StatisticOfReadQuran Struct
type StatisticOfReadQuran struct {
	Type      string `json:"type"`
	Count     int    `json:"count"`
	GroupID   int    `json:"group_id"`
	UserID    int    `json:"user_id"`
	MaxMember int    `json:"max_member"`
	Khatam    bool   `json:"khatam"`
}

// ResponseAllQuran Struct
type ResponseAllQuran struct {
	ID           int
	Asma         string
	SurahName    string
	Transalation string
	Ayat         int
	Number       int
}

// ResponseTotalKhatam Struct
type ResponseTotalKhatam struct {
	Count     int    `json:"count"`
	GroupType string `json:"type"`
	IsJoined  bool   `json:"is_joined"`
}

// ResponseUserTotalKhatam struct
type ResponseUserTotalKhatam struct {
	GroupID   int    `json:"group_id"`
	GroupType string `json:"group_type"`
	Count     int    `json:"count"`
}

// ResponsePersonalReadStatus Struct
type ResponsePersonalReadStatus struct {
	Count int    `json:"count"`
	Type  string `json:"type"`
}

// ResponseListIbuKota struct
type ResponseListIbuKota struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ResponseListProvinces struct
type ResponseListProvinces struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	InternationID string `json:"internation_id"`
}

// ResponseListCities
type ResponseListCities struct {
	ID         int    `json:"id"`
	ProvinceID int    `json:"province_id"`
	Name       string `json:"name"`
}

// ResponseGetKotaFromAPIBangsa struct
type ResponseGetKotaFromAPIBangsa struct {
	Status string
	Query  Query
	Kota   []Kota
}

// Query struct
type Query struct {
	Format string `json:"format"`
	Name   string `json:"name"`
}

// Kota struct
type Kota struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// RequestIDCityByProvinceID struct
type RequestIDCityByProvinceID struct {
	ID int `json:"province_id"`
}

// ResponseTotalUserByProvince struct
type ResponseTotalUserByProvince struct {
	Count         int    `json:"count"`
	InternationID string `json:"internation_id"`
	Name          string `json:"name"`
}

// BoardingPage struct
type BoardingPage struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Position    int       `json:"position"`
	IsActive    int       `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RequestIDBoardingPage struct
type RequestIDBoardingPage struct {
	ID int `json:"id"`
}
