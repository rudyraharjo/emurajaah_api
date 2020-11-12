package service

import (
	"time"

	"github.com/rudyraharjo/emurojaah/admin/province"
	"github.com/rudyraharjo/emurojaah/models"
)

type provinceService struct {
	provinceRepo province.Repository
}

// NewCountryService func
func NewProvinceService(repoProvince province.Repository) province.Service {
	return &provinceService{repoProvince}
}

func (s *provinceService) GetProvinceList() []models.Province {

	res, err := s.provinceRepo.GetProvinceListAll()
	if err != nil {
		return nil
	}
	return res
}

func (s *provinceService) AddProvince(province models.Province) ([]models.Province, error) {

	now := time.Now()

	GetLastProvince, errLastProvince := s.provinceRepo.GetLastProvinceID()

	if errLastProvince != nil {
		return GetLastProvince, errLastProvince
	}

	addProvince := models.Province{
		ID:            GetLastProvince[0].ID + 1,
		Name:          province.Name,
		InternationID: province.InternationID,
		CountryID:     province.CountryID,
		CreatedAt:     now,
	}

	data, err := s.provinceRepo.AddProvince(addProvince)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *provinceService) UpdateProvince(province models.Province) ([]models.Province, error) {

	now := time.Now()

	updateProvince := models.Province{
		ID:            province.ID,
		Name:          province.Name,
		InternationID: province.InternationID,
		CountryID:     province.CountryID,
		UpdatedAt:     now,
	}

	data, err := s.provinceRepo.UpdateProvince(updateProvince)

	if err != nil {
		return data, err
	}
	return data, nil

}

func (s *provinceService) DeletedProvince(paramID int) error {

	err := s.provinceRepo.DeletedProvince(paramID)
	if err != nil {
		return err
	}

	return err

}

func (s *provinceService) GetProvinceByID(paramID int) ([]models.Province, error) {
	data, err := s.provinceRepo.GetProvinceByID(paramID)
	if err != nil {
		return data, err
	}

	return data, nil
}
