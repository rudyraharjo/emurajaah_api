package city

import "github.com/rudyraharjo/emurojaah/models"

// Repository interface
type Repository interface {
	GetCityByID(paramID int) ([]models.City, error)
	GetCityListAll() ([]models.City, error)
	DeletedCity(paramID int) error
	AddCity(country models.City) ([]models.City, error)
	UpdateCity(country models.City) ([]models.City, error)
	GetLastCityID() ([]models.City, error)
}
