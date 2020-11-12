package point

import (
	"github.com/rudyraharjo/emurojaah/models"
)

// Repository interface
type Repository interface {
	GetPointList() ([]models.ResponsePoint, error)
	AddPoint(params models.ResponsePoint) ([]models.ResponsePoint, error)
	UpdateToinactive(paramID int) error
	UpdateToActive(paramID int) error
	UpdatePoint(params models.ResponsePoint) ([]models.ResponsePoint, error)
}
