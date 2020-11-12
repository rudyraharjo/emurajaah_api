package service

import (
	"fmt"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/country"
	"github.com/rudyraharjo/emurojaah/models"
)

type countryService struct {
	countryRepo country.Repository
}

// NewCountryService func
func NewCountryService(repoCountry country.Repository) country.Service {
	return &countryService{repoCountry}
}

func (s *countryService) GetCountryList() []models.Country {

	res, err := s.countryRepo.GetCountryListAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (s *countryService) AddCountry(country models.Country) ([]models.Country, error) {

	now := time.Now()

	addCountry := models.Country{
		CountryName: country.CountryName,
		CountryCode: country.CountryCode,
		CreatedAt:   now,
	}

	data, err := s.countryRepo.AddCountry(addCountry)
	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

func (s *countryService) UpdateCountry(country models.Country) ([]models.Country, error) {

	now := time.Now()

	updateCountry := models.Country{
		ID:          country.ID,
		CountryName: country.CountryName,
		CountryCode: country.CountryCode,
		UpdatedAt:   now,
	}

	data, err := s.countryRepo.UpdateCountry(updateCountry)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (s *countryService) DeletedCountry(paramID int) error {

	err := s.countryRepo.DeletedCountry(paramID)
	if err != nil {
		return err
	}

	return err

}

func (s *countryService) GetCountryByID(paramID int) ([]models.Country, error) {
	data, err := s.countryRepo.GetCountryByID(paramID)
	if err != nil {
		return data, err
	}

	return data, nil
}
