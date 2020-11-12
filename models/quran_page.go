package models

type ResponseFromApiPageQuran struct {
	Code   int      `json:"code"`
	Status string   `json:"status"`
	Data   RAPQData `json:"data"`
}

type RAPQData struct {
	Number int            `json:"number"`
	Ayahs  []RAPQDataAyah `json:"ayahs"`
}

type RAPQDataAyah struct {
	Number        int               `json:"number"`
	Text          string            `json:"text"`
	Surah         RAPQDataAyahSurah `json:"surah"`
	NumberInSurah int               `json:"numberInSurah"`
	Juz           int               `json:"juz"`
	Manzil        int               `json:"manzil"`
	Page          int               `json:"page"`
	Ruku          int               `json:"ruku"`
	HizbQuarter   int               `json:"hizbQuarter"`
}

type RAPQDataAyahSurah struct {
	Number                 int    `json:"number"`
	Name                   string `json:"name"`
	EnglishName            string `json:"englishName"`
	EnglishNameTranslation string `json:"englishNameTranslation"`
	RevelationType         string `json:"revelationType"`
	NumberOfAyahs          int    `json:"numberOfAyahs"`
}
