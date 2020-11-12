package banner

import (
	"github.com/rudyraharjo/emurojaah/models"
)

// Service interface
type Service interface {
	GetBannerList() []models.Banner

	AddBanner(params models.Banner) ([]models.Banner, error)

	UpdateBanner(params models.Banner) ([]models.Banner, error)

	UpdateToinactive(paramID int) error

	UpdateToActive(paramID int) error

	DeleteBanner(paramID int) error
}
