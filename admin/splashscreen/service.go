package splashscreen

import "github.com/rudyraharjo/emurojaah/models"

// Service interface
type Service interface {
	GetSplashScreenList() []models.SplashScreen
	GetSplashScreenListAll() []models.SplashScreen
	AddSplashScreen(param models.SplashScreen) ([]models.SplashScreen, error)
	UpdateSplashScreen(splash models.SplashScreen) ([]models.SplashScreen, error)
	UpdateToinactive(paramID int) error
	UpdateToActive(paramID int) error
	DeleteSplashScreen(id int) error
}
