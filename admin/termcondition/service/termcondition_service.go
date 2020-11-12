package service

import (
	"time"

	"github.com/rudyraharjo/emurojaah/admin/termcondition"
	"github.com/rudyraharjo/emurojaah/models"
)

type termconditionService struct {
	termconditionRepo termcondition.Repository
}

// NewTermConditionService func
func NewTermConditionService(repoTermCondition termcondition.Repository) termcondition.Service {
	return &termconditionService{repoTermCondition}
}

func (s *termconditionService) GetListTermConditions() []models.TermCondition {

	res, err := s.termconditionRepo.GetListTermConditions()
	if err != nil {
		return nil
	}
	return res

}

func (s *termconditionService) Update(term models.TermCondition) ([]models.TermCondition, error) {

	now := time.Now()

	update := models.TermCondition{
		ID:          term.ID,
		Description: term.Description,
		UpdatedAt:   now,
	}

	data, err := s.termconditionRepo.Update(update)

	if err != nil {
		return data, err
	}
	return data, nil

}
