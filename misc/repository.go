package misc

import "github.com/rudyraharjo/emurojaah/models"

// Repository get list province
type Repository interface {
	GetProvinceList() ([]models.Province, error)
}
