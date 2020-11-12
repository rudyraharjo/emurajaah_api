package banner

import (
	"github.com/rudyraharjo/emurojaah/models"
)

// Repository interface
type Repository interface {
	GetBannerList() ([]models.Banner, error)
	AddBanner(params models.Banner) ([]models.Banner, error)
	UpdateBanner(params models.Banner) ([]models.Banner, error)
	UpdateToinactive(paramID int) error
	UpdateToActive(paramID int) error
	DeleteBanner(paramID int) error

	// ----- noroute ----- //
	GetBannerListActive() ([]models.Banner, error)
}
