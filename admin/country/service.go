package country

import "github.com/rudyraharjo/emurojaah/models"

// Service interface
type Service interface {
	GetCountryByID(paramID int) ([]models.Country, error)
	GetCountryList() []models.Country
	AddCountry(param models.Country) ([]models.Country, error)
	DeletedCountry(paramID int) error
	UpdateCountry(param models.Country) ([]models.Country, error)
}
