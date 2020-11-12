package province

import "github.com/rudyraharjo/emurojaah/models"

// Repository interface
type Repository interface {
	GetProvinceByID(paramID int) ([]models.Province, error)
	GetProvinceListAll() ([]models.Province, error)
	DeletedProvince(paramID int) error
	AddProvince(country models.Province) ([]models.Province, error)
	UpdateProvince(country models.Province) ([]models.Province, error)
	GetLastProvinceID() ([]models.Province, error)
}
