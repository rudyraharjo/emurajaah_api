package content

import "github.com/rudyraharjo/emurojaah/models"

type Repository interface {

	// create
	AddAyatQuran(surah []models.AyatQuran) error

	SaveUserReadingActivity(params models.RequestSaveReadingQuran) ([]models.UserActivity, error)

	// get list

	GetSurahByIDSurah(SurahID int) ([]models.ResponseQuran, error)
	GetAyatByCategory(category string, index int) ([]models.ResponseQuran, error)
	GetQuranByCategoryAndPaging(category string, index int, offset int, limit int) ([]models.ResponseQuran, error)
	GetQuotes() ([]models.Quote, error)
	ListUserGroup(userId int) ([]models.ResponseGroupList, error)

	// update
	UpdateQuranPageById(id int, page int) error

	// validation

	// GetAllQuranList
	GetAllQuranList() ([]models.ResponseAllQuran, error)

	// --- SplashScreen ---

	//GetSplashScreenList
	GetSplashScreenList() ([]models.SplashScreen, error)

	// GetTotalKhatamAllGroup
	GetTotalKhatamAllGroup() ([]models.ResponseTotalKhatam, error)

	// GetListIbukota
	GetListIbukota() ([]models.ResponseListIbuKota, error)

	//GetListProvinces
	GetListProvinces() ([]models.ResponseListProvinces, error)

	// GetListCities
	GetListCities(provinceID int) ([]models.ResponseListCities, error)

	// GetTotalUserByProvince
	GetTotalUserByProvince() ([]models.ResponseTotalUserByProvince, bool)

	// GetListAllBoarding
	GetListAllBoarding() ([]models.BoardingPage, error)

	// GetListBoardingPage
	GetListBoardingPageIsActive() ([]models.BoardingPage, error)

	// AddBoardingPage
	AddBoardingPage(param models.BoardingPage) ([]models.BoardingPage, error)

	// UpdateBoardingPage
	UpdateBoardingPage(param models.BoardingPage) ([]models.BoardingPage, error)

	// DeactivateBoardingPage
	DeactivateBoardingPage(paramID int) error
}
