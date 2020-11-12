package models

import "time"

type SurahRequest struct {
	Ar    string `json:"ar"`
	Id    string `json:"id"`
	Nomor string `json:"nomor"`
	Tr    string `json:"tr"`
}

type AddQuranRequest struct {
	AyatId int            `json:"ayat_id"`
	Surah  []SurahRequest `json:"surah"`
}

type RequestQuran struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type RequestQuranPaging struct {
	Type   string `json:"type"`
	Index  int    `json:"index"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type ResponseQuran struct {
	Arabic      string `json:"arabic"`
	Latin       string `json:"latin"`
	Translation string `json:"translation"`
	Audio       string `json:"audio"`
	Image       string `json:"image"`
	AyatSec     int    `json:"sequence"`
	SurahName   string `json:"surah_name"`
	SurahType   string `json:"surah_type"`
	Number      int    `json:"surah_number"`
	Ayat        int    `json:"total_ayat"`
}

type AyatQuran struct {
	Id          int
	Arabic      string
	Latin       string
	Translation string
	Audio       string
	Image       string
	JuzId       int
	SurahId     int
	AyatSec     int
	Page        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RequestQuranBySurahId Struct
type RequestQuranBySurahId struct {
	SurahID int `json:"surah_id"`
}
