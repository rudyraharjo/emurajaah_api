package service

import (
	"fmt"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/splashscreen"
	"github.com/rudyraharjo/emurojaah/models"
)

type splashscreenService struct {
	splashscreenRepo splashscreen.Repository
}

// NewSplashScreenService func
func NewSplashScreenService(repoSplashScreen splashscreen.Repository) splashscreen.Service {
	return &splashscreenService{repoSplashScreen}
}

func (s *splashscreenService) GetSplashScreenListAll() []models.SplashScreen {

	res, err := s.splashscreenRepo.GetSplashScreenListAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res

}

func (s *splashscreenService) GetSplashScreenList() []models.SplashScreen {

	res, err := s.splashscreenRepo.GetSplashScreenList()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (s *splashscreenService) AddSplashScreen(splashscreen models.SplashScreen) ([]models.SplashScreen, error) {

	now := time.Now()

	addSplashScreen := models.SplashScreen{
		Title:       splashscreen.Title,
		Description: splashscreen.Description,
		ImageURL:    splashscreen.ImageURL,
		Position:    splashscreen.Position,
		IsActive:    1,
		CreatedAt:   now,
	}

	data, err := s.splashscreenRepo.AddSplashScreen(addSplashScreen)
	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

func (s *splashscreenService) UpdateSplashScreen(splashscreen models.SplashScreen) ([]models.SplashScreen, error) {

	now := time.Now()

	SplashScreen := models.SplashScreen{
		ID:          splashscreen.ID,
		Title:       splashscreen.Title,
		Description: splashscreen.Description,
		ImageURL:    splashscreen.ImageURL,
		Position:    splashscreen.Position,
		IsActive:    1,
		UpdatedAt:   now,
	}

	data, err := s.splashscreenRepo.UpdateSplashScreen(SplashScreen)

	if err != nil {
		fmt.Println(err)
		return data, err
	}
	return data, nil

}

func (s *splashscreenService) DeleteSplashScreen(id int) error {
	err := s.splashscreenRepo.DeleteSplashScreen(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *splashscreenService) UpdateToinactive(paramID int) error {

	err := s.splashscreenRepo.UpdateToinactive(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (s *splashscreenService) UpdateToActive(paramID int) error {

	err := s.splashscreenRepo.UpdateToActive(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
