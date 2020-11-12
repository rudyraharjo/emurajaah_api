package point

import (
	"github.com/rudyraharjo/emurojaah/models"
)

// Service interface
type Service interface {
	GetPointList() []models.ResponsePoint

	AddPoint(params models.RequestAddPoint) ([]models.ResponsePoint, error)

	UpdateToinactive(paramID int) error

	UpdateToActive(paramID int) error

	UpdatePoint(params models.ResponsePoint) ([]models.ResponsePoint, error)
}
