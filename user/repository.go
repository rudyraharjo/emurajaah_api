package user

import "github.com/rudyraharjo/emurojaah/models"

type Repository interface {
	//Get List
	RetrieveCredentialAdminByAlias(alias string) (*models.UserCredential, error)
	RetrieveCredentialByAlias(alias string) (*models.UserCredential, error)
	GetUserIdByAlias(alias string) (int, error)
	GetUserBasicInfoById(userId int) (*models.ResponseUserBasicInfo, error)
	GetUserPointGroupByType(userId int) ([]models.UserReward, error)
	GetUserAliasById(userId int) ([]models.UserAlias, error)
	GetPublicStatOfReadQuran() ([]models.ResponseTotalKhatam, error)
	GetPublicStatOfReadQuranByUserID(userID int) ([]models.ResponseTotalKhatam, error)
	GetUserStatOfReadQuran(userId int) ([]models.ResponsePersonalReadStatus, error)

	//Validation
	IsUserExist(id int) (int, error)
	IsAliasExist(alias string) (int, error)

	//Create
	AddUser(user models.UserBasicInfo) (int, error)
	AddUserAlias(alias []models.UserAliasFull) error
	AddSingleUserAlias(alias models.UserAliasFull) error
	AddUserToken(alias models.UserToken) error
	AddSurvey(survey models.UserSurvey) error

	//Update
	UpdateProfile(userInfo models.UserBasicInfo) error

	//checkandaddToken
	IsTokenExits(id int) (int, error)

	//DeleteTokenFirebase
	DeleteTokenFirebase(idUser int) error

	//CheckSurahIsDone
	CheckSurahIsDone(params models.ReadIsDone) (int, error)

	// GetUserAdminByAlias
	GetListUserAdminByAlias() ([]models.UserAdminAliasFull, error)
	GetUserAdminByAlias(alias string) (*models.UserAdminAliasFull, error)

	//GetUserMemberList
	GetUserMemberList() ([]models.UserMemberAliasFull, error)

	// GetUserPointGroupByType
	GetPointGroupByType() ([]models.UserReward, error)

	//GetTotalMember
	GetTotalMember() (int, error)

	GetSurveyJuzzGrouping() []models.ResponseSurveyJuzzGrouping

	GetNameProvinceFromUser() ([]models.User, error)
}
