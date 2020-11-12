package misc

import "github.com/rudyraharjo/emurojaah/models"

// Service get list province
type Service interface {
	GetProvinceList() []models.Province
}
