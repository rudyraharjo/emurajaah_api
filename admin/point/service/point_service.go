package service

import (
	"fmt"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/point"
	"github.com/rudyraharjo/emurojaah/models"
)

type pointService struct {
	pointRepo point.Repository
}

// NewPointService func
func NewPointService(repo point.Repository) point.Service {
	return &pointService{repo}
}

func (s *pointService) GetPointList() []models.ResponsePoint {
	data, err := s.pointRepo.GetPointList()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *pointService) AddPoint(params models.RequestAddPoint) ([]models.ResponsePoint, error) {

	now := time.Now()

	PointAdd := models.ResponsePoint{
		Type:      params.Type,
		Point:     params.Point,
		CreatedAt: now,
		IsActive:  1,
	}

	data, err := s.pointRepo.AddPoint(PointAdd)

	return data, err
}

func (s *pointService) UpdatePoint(params models.ResponsePoint) ([]models.ResponsePoint, error) {
	now := time.Now()

	PointUpdate := models.ResponsePoint{
		ID:        params.ID,
		Type:      params.Type,
		Point:     params.Point,
		UpdatedAt: now,
		IsActive:  1,
	}

	data, err := s.pointRepo.UpdatePoint(PointUpdate)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *pointService) UpdateToinactive(paramID int) error {

	err := s.pointRepo.UpdateToinactive(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (s *pointService) UpdateToActive(paramID int) error {

	err := s.pointRepo.UpdateToActive(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
