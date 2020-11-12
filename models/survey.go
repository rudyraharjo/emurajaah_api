package models

import "time"

type UserSurvey struct {
	Id             int
	UserId         int
	SurveyType     string
	SurveyQuestion string
	SurveyAnswer   string
	SurveyImage    string
	CreatedDate    time.Time
}

type RequestJuzzSurvey struct {
	UserId int    `json:"user_id"`
	Answer string `json:"answer"`
}
