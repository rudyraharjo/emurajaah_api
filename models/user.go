package models

import "time"

type User struct {
	ID             int
	FullName       string
	Gender         string
	BirthDate      time.Time
	BirthPlace     string
	Address        string
	Province       string
	ProfilePicture string
	Alias          []UserAlias
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserAlias struct {
	AliasType string
	Alias     string
}

// UserBasicInfo struct
type UserBasicInfo struct {
	Id             int       `json:"id"`
	FullName       string    `json:"full_name"`
	Gender         string    `json:"gender"`
	BirthDate      time.Time `json:"birth_date"`
	BirthPlace     string    `json:"birth_place"`
	Address        string    `json:"address"`
	IDProvince     int       `json:"id_province"`
	IDCity         int       `json:"id_city"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ResponseUserBasicInfo struct
type ResponseUserBasicInfo struct {
	ID             int       `json:"id"`
	FullName       string    `json:"full_name"`
	Gender         string    `json:"gender"`
	BirthDate      time.Time `json:"birth_date"`
	BirthPlace     string    `json:"birth_place"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	Province       string    `json:"province"`
	IDProvince     int       `json:"id_province"`
	IDCity         int       `json:"id_city"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserBasicInfoWithEmail struct {
	Id             int       `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Gender         string    `json:"gender"`
	BirthDate      time.Time `json:"birth_date"`
	BirthPlace     string    `json:"birth_place"`
	Address        string    `json:"address"`
	Province       string    `json:"province"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserAliasFull struct {
	UserId         int
	AliasType      string
	Alias          string
	Credential     string
	CredentialType string
	IsVerified     int
	ActivationCode string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// UserAdminAliasFull Struct
type UserAdminAliasFull struct {
	UserID          int
	Username        string
	Alias           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ActivatedStatus bool
	RoleName        string
}

//UserMemberAliasFull Struct
type UserMemberAliasFull struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	Cityname  string    `json:"cityname"`
	Provname  string    `json:"provname"`
	CreatedAt time.Time `json:"created_at"`
}

// UserToken Type struct
type UserToken struct {
	UserID        int
	TokenFirebase string
	CreatedAt     time.Time
}

type UserCredential struct {
	UserID     int
	Credential string
}

// UserActivity struct
type UserActivity struct {
	Id               int       `json:"id"`
	IdGroupMember    int       `json:"id_group_member"`
	UserId           int       `json:"user_id"`
	GroupId          int       `json:"group_id"`
	GroupType        string    `json:"group_type"`
	ContentIndex     int       `json:"content_index"`
	StatusUserAction int       `json:"status_user_action"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserReward struct {
	Type      string `json:"type"`
	Point     int    `json:"point"`
	UserPoint int    `json:"user_point"`
}

type RequestByUserId struct {
	UserId        int    `json:"user_id"`
	TokenFirebase string `json:"token_firebase"`
}

// ReadIsDone type struct
type ReadIsDone struct {
	UserID  int `json:"user_id"`
	GroupID int `json:"group_id"`
	Index   int `json:"index"`
}

// RequestLogin type struct
type RequestLogin struct {
	Alias         string `json:"alias"`
	Password      string `json:"password"`
	Type          string `json:"type"`
	TokenFirebase string `json:"token_firebase"`
}

// RequestRegister struct
type RequestRegister struct {
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Gender        string `json:"gender"`
	BirthDate     string `json:"birth_date"`
	Address       string `json:"address"`
	IDProvince    int    `json:"id_province"`
	IDCity        int    `json:"id_city"`
	TokenFirebase string `json:"token_firebase"`
}

type RequestGetProfile struct {
	Alias string `json:"alias"`
}

// RequestEditProfile struct
type RequestEditProfile struct {
	UserId         int    `json:"user_id"`
	FullName       string `json:"full_name"`
	Gender         string `json:"gender"`
	BirthDate      string `json:"birth_date"`
	BirthPlace     string `json:"birth_place"`
	Address        string `json:"address"`
	IDProvince     int    `json:"id_province"`
	IDCity         int    `json:"id_city"`
	ProfilePicture string `json:"profile_picture"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}

// ResponseGetProfileOverview struct
type ResponseGetProfileOverview struct {
	Profile                ResponseUserBasicInfo        `json:"profile"`
	PersonalReadStatus     []ResponsePersonalReadStatus `json:"personal_read_status"`
	GlobalGroupTotalKhatam []ResponseTotalKhatam        `json:"global_group_status"`
}

// ResponseListAdmin Struct
type ResponseListAdmin struct {
	UserListAdmin []UserAdminAliasFull
}

// ResponseUserAdmin Struct
type ResponseUserAdmin struct {
	FullName string
	gender   string
	address  string
	Province string
	UserName string
	Alias    string
	RoleName string
}

// UserRewardAyat Struct
type UserRewardAyat struct {
	Type      string `json:"type"`
	Point     int    `json:"point"`
	UserPoint int    `json:"user_point"`
}

// UserRewardHalaman Struct
type UserRewardHalaman struct {
	Type      string `json:"type"`
	Point     int    `json:"point"`
	UserPoint int    `json:"user_point"`
}

// UserRewardJuz Struct
type UserRewardJuz struct {
	Type      string `json:"type"`
	Point     int    `json:"point"`
	UserPoint int    `json:"user_point"`
}

// ResponseSurveyJuzzGrouping struct
type ResponseSurveyJuzzGrouping struct {
	JuzID    int `json:"juz_id"`
	CountJuz int `json:"count"`
}

// ResponseSurveyJuz struct
type ResponseSurveyJuz struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	SurveyType     string    `json:"survey_type"`
	SurveyQuestion string    `json:"survey_question"`
	SurveyAnswer   string    `json:"survey_answer"`
	SurveyImage    string    `json:"survey_image"`
	CreatedAt      time.Time `json:"created_at"`
}
