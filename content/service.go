package content

import "github.com/rudyraharjo/emurojaah/models"

type Service interface {

	// main control
	HandlerAutoUpdateQuranPage()

	HandlerUpdateIDCityAndIDProv()

	// ForHome
	HandlerHomePageContentBanner(userID int) []models.Banner
	HandlerHomePageContentGroups(userID int) []models.ResponseGroupUserJoined
	HandlerHomePageContentQuotes(userID int) []models.Quote
	HandlerHomePageContentGlobalGroupStatus(userID int) []models.ResponseTotalKhatam

	HandlerHomePageContent(userId int) models.ResponseHomePageContent

	// create
	AddQuranSurah(reqParams models.AddQuranRequest) error

	SaveUserReadingActivity(params models.RequestSaveReadingQuran) error

	// get list

	GetSurahByIDSurah(surahID int) []models.ResponseQuran
	GetSurahByCategory(params models.RequestQuran) []models.ResponseQuran
	GetQuranByCategoryWithPaging(params models.RequestQuranPaging) []models.ResponseQuran
	GetQuoteList() []models.Quote
	GetListUserGroup(userId int) []models.ResponseGroupList

	// Update
	UpdateQuranPageById(id int, page int) error

	//UpdateGroupMemberReadingStatus(groupId int, userId int, status int) error

	// validation

	// http request
	GetQuranFromAPI() error
	GetQuranPageFromAPI(index int) error

	// getAllQuran
	GetAllQuran() []models.ResponseAllQuran

	// TERM-CONDITION
	GetTermCondition() []models.TermCondition
	// ---- SplashScreen ----

	//GetSplashScreenList
	GetSplashScreenList() []models.SplashScreen

	// GetTotalKhatamAllGroup
	GetTotalKhatamAllGroup() ([]models.ResponseTotalKhatam, error)

	//GetListIbukota
	GetListIbukota() []models.ResponseListIbuKota

	//GetListProvinces
	GetListProvinces() []models.ResponseListProvinces

	// GetListCities
	GetListCities(provinceID int) []models.ResponseListCities

	// GetTotalUserByProvince
	GetTotalUserByProvince() ([]models.ResponseTotalUserByProvince, bool)

	// GetListAllBoarding
	GetListAllBoarding() []models.BoardingPage

	// GetListBoardingPage
	GetListBoardingPageIsActive() []models.BoardingPage

	// AddBoardingPage
	AddBoardingPage(param models.BoardingPage) ([]models.BoardingPage, error)

	// UpdateBoardingPage
	UpdateBoardingPage(params models.BoardingPage) ([]models.BoardingPage, error)

	// DeactivateBoardingPage
	DeactivateBoardingPage(boardingpageID int) error
}
