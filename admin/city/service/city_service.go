package service

import (
	"time"

	"github.com/rudyraharjo/emurojaah/admin/city"
	"github.com/rudyraharjo/emurojaah/models"
)

type cityService struct {
	cityRepo city.Repository
}

// NewCityService func
func NewCityService(repoCity city.Repository) city.Service {
	return &cityService{repoCity}
}

func (s *cityService) GetCityList() []models.City {

	res, err := s.cityRepo.GetCityListAll()
	if err != nil {
		return nil
	}
	return res
}

func (s *cityService) AddCity(city models.City) ([]models.City, error) {

	now := time.Now()

	GetLastCity, errLastCity := s.cityRepo.GetLastCityID()

	if errLastCity != nil {
		return GetLastCity, errLastCity
	}

	AddCity := models.City{
		ID:         GetLastCity[0].ID + 1,
		ProvinceID: city.ProvinceID,
		Name:       city.Name,
		CreatedAt:  now,
	}

	data, err := s.cityRepo.AddCity(AddCity)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *cityService) UpdateCity(city models.City) ([]models.City, error) {

	now := time.Now()

	updateCity := models.City{
		ID:         city.ID,
		ProvinceID: city.ProvinceID,
		Name:       city.Name,
		UpdatedAt:  now,
	}

	data, err := s.cityRepo.UpdateCity(updateCity)

	if err != nil {
		return data, err
	}
	return data, nil

}

func (s *cityService) DeletedCity(paramID int) error {

	err := s.cityRepo.DeletedCity(paramID)
	if err != nil {
		return err
	}

	return err

}

func (s *cityService) GetCityByID(paramID int) ([]models.City, error) {
	data, err := s.cityRepo.GetCityByID(paramID)
	if err != nil {
		return data, err
	}

	return data, nil
}
