package termcondition

import "github.com/rudyraharjo/emurojaah/models"

// Repository interface
type Repository interface {
	GetListTermConditionsIsActived() ([]models.TermCondition, error)
	GetListTermConditions() ([]models.TermCondition, error)
	Update(termParams models.TermCondition) ([]models.TermCondition, error)
}
