package city

import "github.com/rudyraharjo/emurojaah/models"

// Service interface
type Service interface {
	GetCityByID(paramID int) ([]models.City, error)
	GetCityList() []models.City
	AddCity(param models.City) ([]models.City, error)
	DeletedCity(paramID int) error
	UpdateCity(param models.City) ([]models.City, error)
}
