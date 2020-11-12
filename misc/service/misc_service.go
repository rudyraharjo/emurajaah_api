package service

import (
	"fmt"

	"github.com/rudyraharjo/emurojaah/misc"
	"github.com/rudyraharjo/emurojaah/models"
)

type miscService struct {
	miscRepo misc.Repository
}

func NewMiscService(repo misc.Repository) misc.Service {
	return &miscService{repo}
}

func (s *miscService) GetProvinceList() []models.Province {
	data, err := s.miscRepo.GetProvinceList()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}
