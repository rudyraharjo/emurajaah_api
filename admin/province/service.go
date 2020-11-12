package province

import "github.com/rudyraharjo/emurojaah/models"

// Service interface
type Service interface {
	GetProvinceByID(paramID int) ([]models.Province, error)
	GetProvinceList() []models.Province
	AddProvince(param models.Province) ([]models.Province, error)
	DeletedProvince(paramID int) error
	UpdateProvince(param models.Province) ([]models.Province, error)
}
