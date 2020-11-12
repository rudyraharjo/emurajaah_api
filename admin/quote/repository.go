package quote

import "github.com/rudyraharjo/emurojaah/models"

// Repository interface
type Repository interface {
	GetQuoteByID(paramID int) ([]models.Quote, error)
	GetQuotesListAll() ([]models.Quote, error)
	DeletedQuote(paramID int) error
	AddQuote(quote models.Quote) ([]models.Quote, error)
	UpdateQuote(quote models.Quote) ([]models.Quote, error)
	UpdateStatus(paramID int) error
	// GetSplashScreenListAll() ([]models.SplashScreen, error)
	// GetSplashScreenListActive() ([]models.SplashScreen, error)
	// AddSplashScreen(splash models.SplashScreen) ([]models.SplashScreen, error)
	// UpdateSplashScreen(splash models.SplashScreen) ([]models.SplashScreen, error)
	// UpdateToinactive(paramID int) error
	// UpdateToActive(paramID int) error
	// DeleteSplashScreen(id int) error
}
