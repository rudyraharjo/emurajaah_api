package country

import "github.com/rudyraharjo/emurojaah/models"

// Repository interface
type Repository interface {
	GetCountryByID(paramID int) ([]models.Country, error)
	GetCountryListAll() ([]models.Country, error)
	DeletedCountry(paramID int) error
	AddCountry(country models.Country) ([]models.Country, error)
	UpdateCountry(country models.Country) ([]models.Country, error)
}
