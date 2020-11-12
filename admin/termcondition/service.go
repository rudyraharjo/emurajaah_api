package termcondition

import "github.com/rudyraharjo/emurojaah/models"

// Service interface
type Service interface {
	GetListTermConditions() []models.TermCondition
	Update(params models.TermCondition) ([]models.TermCondition, error)
	// GetSplashScreenListAll() []models.SplashScreen
	// AddSplashScreen(param models.SplashScreen) ([]models.SplashScreen, error)
	// UpdateToinactive(paramID int) error
	// UpdateToActive(paramID int) error
	// DeleteSplashScreen(id int) error
}
