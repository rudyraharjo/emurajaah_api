package user

import "github.com/rudyraharjo/emurojaah/models"

type Service interface {

	//Main Service
	RegisterUser(registerParams models.RequestRegister) (int, error)
	AddJuzzSurvey(reqParams models.RequestJuzzSurvey) error
	LoginUser(params models.RequestLogin) (bool, int)
	LoginUserWithGoogle(params models.RequestLogin) (bool, int)
	LoginUserAdmin(params models.RequestLogin) (bool, int)

	//Validation
	CheckIsUserExist(id int) int
	CheckIsAliasExist(alias string) int
	CheckPasswordHash(password, hash string) bool
	CheckTokenFirebase(id int) int

	//Get List
	GetUserBasicInfoByAlias(alias string) (*models.ResponseUserBasicInfo, error)
	GetUserBasicInfoById(userId int) (*models.ResponseUserBasicInfo, error)
	GetUserPointReward(userId int) (int, []models.UserReward, error)
	GetUserAlias(userId int) []models.UserAlias
	GetPublicGroupStatistic(userID int) []models.ResponseTotalKhatam
	GetUserReadStatistic(userId int) []models.ResponsePersonalReadStatus

	//Update
	UpdateUserProfile(userInfo models.RequestEditProfile) error

	// other
	HashPassword(password string) (string, error)
	LoginWithRequest(params models.RequestLogin) (int, *models.ResponseLoginSuccess, *models.ResponseLoginFail)

	//AddTokenFirebase
	AddTokenUser(id int, token string) (int, error)

	//DeleteTokenFirebase
	DeleteTokenFirebase(idUser int) error

	//CheckSurahIsDone
	CheckSurahReadingIsDone(readisdoneParams models.ReadIsDone) (int, error)

	//  --- Admin ---

	// GetUserAdminInfoByAlias
	GetUserAdminInfoByAlias(alias string) (*models.UserAdminAliasFull, *models.ResponseUserBasicInfo, error)

	// GetUserAdminList
	GetUserAdminList() []models.UserAdminAliasFull

	//GetUserAdminList
	GetUserMemberList() []models.UserMemberAliasFull

	// GetUserPointReward
	GetPointReward() (int, []models.UserReward, error)

	// GetTotalMember
	GetTotalMember() (int, error)

	// Survey Juz
	SurveyJuzzGrouping() ([]models.ResponseSurveyJuzzGrouping, error)
}
