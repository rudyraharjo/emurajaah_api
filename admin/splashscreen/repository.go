package splashscreen

import "github.com/rudyraharjo/emurojaah/models"

// Repository interface
type Repository interface {
	GetSplashScreenList() ([]models.SplashScreen, error)
	GetSplashScreenListAll() ([]models.SplashScreen, error)
	GetSplashScreenListActive() ([]models.SplashScreen, error)
	AddSplashScreen(splash models.SplashScreen) ([]models.SplashScreen, error)
	UpdateSplashScreen(splash models.SplashScreen) ([]models.SplashScreen, error)
	UpdateToinactive(paramID int) error
	UpdateToActive(paramID int) error
	DeleteSplashScreen(id int) error
}
