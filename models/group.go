package models

import "time"

// GroupMember struct
type GroupMember struct {
	ID           int       `json:"id"`
	GroupID      int       `json:"group_id"`
	UserID       int       `json:"user_id"`
	CurrentIndex int       `json:"current_index"`
	GroupType    string    `json:"group_type"`
	IsDone       int       `json:"is_reading"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RequestIDGroupMember struct
type RequestIDGroupMember struct {
	ID int `json:"id_group_member"`
}

// RequestLeaveReading struct
type RequestLeaveReading struct {
	ID           int    `json:"id_group_member"`
	UserID       int    `json:"user_id"`
	GroupID      int    `json:"group_id"`
	GroupType    string `json:"group_type"`
	ContentIndex int    `json:"content_index"`
}

type GroupMemberWithType struct {
	Id           int
	GroupId      int
	GroupType    string
	UserId       int
	CurrentIndex int
	IsDone       int
	IsActive     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GroupMemberCheckRead type struct
type GroupMemberReadNotIsDone struct {
	UserID        int
	IsDone        int
	FullName      string
	TokenFirebase string
}

// Group struct
type Group struct {
	Id            int
	NoGroupIndex  int
	GroupType     string
	GroupName     string
	MaxMember     int
	CurrentMember int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// GroupWithStatus struct
type GroupWithStatus struct {
	ID            int
	NoGroupIndex  int
	GroupType     string
	GroupName     string
	MaxMember     int
	CurrentMember int
	Selesai       int
	Belom         int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// RequestJoinGroupByEmail struct
type RequestJoinGroupByEmail struct {
	GroupID   int    `json:"group_id"`
	Email     string `json:"email"`
	GroupType string `json:"group_type"`
}

// RequestJoinGroup struct
type RequestJoinGroup struct {
	UserID    int    `json:"user_id"`
	GroupType string `json:"group_type"`
}

// RequestListTypeGroupPaging struct
type RequestListTypeGroupPaging struct {
	UserID    int    `json:"user_id"`
	GroupType string `json:"group_type"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// RequestJoinBulkGroup struct
type RequestJoinBulkGroup struct {
	GroupID  int        `json:"id_group"`
	DataBulk []DataBulk `json:"data"`
}

// DataBulk struct
type DataBulk struct {
	Email       string `json:"email"`
}
type RequestListGroup struct {
	UserID int `json:"user_id"`
}

type RequestListMemberGroup struct {
	UserId  int `json:"user_id"`
	GroupId int `json:"group_id"`
}

type RequestListMemberGroupWithOffsetLimit struct {
	UserId  int `json:"user_id"`
	GroupId int `json:"group_id"`
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
}

type RequestListMemberGroupByStatus struct {
	UserId  int `json:"user_id"`
	GroupId int `json:"group_id"`
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Status  int `json:"status"`
}

// RequestSaveReadingQuran struct
type RequestSaveReadingQuran struct {
	ID           int    `json:"id_group_member"`
	UserID       int    `json:"user_id"`
	GroupID      int    `json:"group_id"`
	GroupType    string `json:"group_type"`
	ContentIndex int    `json:"content_index"`
	// UserId    int    `json:"user_id"`
	// GroupId   int    `json:"group_id"`
	// GroupType string `json:"group_type"`
	// Index     int    `json:"index"`
}

type RequestExitGroup struct {
	UserId  int `json:"user_id"`
	GroupId int `json:"group_id"`
}

type ResponseGroupList struct {
	Id            int    `json:"group_id"`
	GroupName     string `json:"group_name"`
	GroupType     string `json:"group_type"`
	CurrentMember int    `json:"members"`
	CurrentIndex  int    `json:"current_index"`
	IsReading     bool   `json:"is_reading"`
}

// ResponsGroupeGroupMembersList struct
type ResponsGroupeGroupMembersList struct {
	UserID    int    `json:"user_id"`
	GroupID   int    `json:"group_id"`
	GroupType string `json:"group_type"`
}

// ResponsGroupeGroupMembersListWithNoUrut struct
type ResponsGroupeGroupMembersListWithNoUrut struct {
	NoUrut    int    `json:"no_urut"`
	UserID    int    `json:"user_id"`
	GroupID   int    `json:"group_id"`
	GroupType string `json:"group_type"`
}

// ResponseListGroupType Struct
type ResponseListGroupType struct {
	ID           int    `json:"id_group_member"`
	GroupID      int    `json:"group_id"`
	UserID       int    `json:"user_id"`
	CurrentIndex int    `json:"current_index"`
	GroupName    string `json:"group_name"`
	GroupType    string `json:"group_type"`
	Asma         string `json:"asma"`
	SurahName    string `json:"surah_name"`
	IsReading    bool   `json:"is_reading"`
}

// ResponseListGroupTypeWithNoUrut struct
type ResponseListGroupTypeWithNoUrut struct {
	NoUrut       int    `json:"no_urut"`
	ID           int    `json:"id_group_member"`
	GroupID      int    `json:"group_id"`
	UserID       int    `json:"user_id"`
	CurrentIndex int    `json:"current_index"`
	GroupName    string `json:"group_name"`
	GroupType    string `json:"group_type"`
	Asma         string `json:"asma"`
	SurahName    string `json:"surah_name"`
	IsReading    bool   `json:"is_reading"`
	JuzID        int    `json:"jus_id"`
	SurahID      int    `json:"surah_id"`
	AyatSec      int    `json:"ayat_sec"`
}

// ResponseAllgroup struct
type ResponseAllgroup struct {
	GroupType string `json:"type"`
	GroupName string `json:"name"`
}

// ResponseGroupUserJoined struct
type ResponseGroupUserJoined struct {
	GroupType string `json:"type"`
	GroupName string `json:"name"`
	IsJoined  bool   `json:"join"`
}

// ResponseGroupMemberList stuct
type ResponseGroupMemberList struct {
	ID             int    `json:"id"`
	UserId         int    `json:"user_id"`
	FullName       string `json:"full_name"`
	Gender         string `json:"gender"`
	ProfilePicture string `json:"profile_picture"`
	CurrentIndex   int    `json:"current_index"`
	IsDone         int    `json:"is_done"`
}

// ResponseGroupMember struct
type ResponseGroupMember struct {
	UserID         int    `json:"user_id"`
	FullName       string `json:"full_name"`
	Gender         string `json:"gender"`
	ProfilePicture string `json:"profile_picture"`
	CurrentIndex   int    `json:"current_index"`
	IsDone         int    `json:"is_done"`
	Email          string `json:"email"`
	JuzID          int    `json:"juz_id"`
	SurahID        int    `json:"surah_id"`
	AyatSec        int    `json:"ayat_sec"`
	Page           int    `json:"page"`
	SurahName      string `json:"surah_name"`
	NoGroupIndex   int    `json:"no_group_index"`
}

type RequestSendFcm struct {
	Notification    NotificationSendFcm `json:"notification"`
	RegistrationIds []string            `json:"registration_ids"`
}

type RequestSingleSendFcm struct {
	Notification NotificationSendFcm `json:"notification"`
	To           string              `json:"to"`
}

type History struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Process     string    `json:"process"`
	CreatedAt   time.Time `json:"created_at"`
}

type NotificationSendFcm struct {
	Body     string `json:"body"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
}

// ResponseGroupListType Struct
type ResponseGroupListType struct {
	GroupType string `json:"group_type"`
}

// GroupKhatam struct
type GroupKhatam struct {
	ID            int       `json:"id"`
	IDRefActivity int       `json:"id_ref_activity"`
	GroupID       int       `json:"group_id"`
	GroupType     string    `json:"group_type"`
	CreatedAt     time.Time `json:"created_at"`
}

// ResponseNotifUserKhatam struct
type ResponseNotifUserKhatam struct {
	NoGroup       int    `json:"no_group"`
	GroupID       int    `json:"group_id"`
	GroupType     string `json:"group_type"`
	Description   string `json:"description"`
	FullName      string `json:"full_name"`
	TokenFirebase string `json:"token_firebase"`
}

// RequestIDGroup struct
type RequestIDGroup struct {
	ID int `json:"group_id"`
}

// RequestTypeGroup struct
type RequestTypeGroup struct {
	GroupType string `json:"group_type"`
}

// GroupMemberMaxMember struct
type GroupMemberMaxMember struct {
	ID           int
	UserID       int
	CurrentIndex int
}
