package service

import (
	"fmt"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/banner"
	"github.com/rudyraharjo/emurojaah/models"
)

type bannerService struct {
	bannerRepo banner.Repository
}

// NewPointService func
func NewBannerService(repo banner.Repository) banner.Service {
	return &bannerService{repo}
}

func (s *bannerService) GetBannerList() []models.Banner {
	data, err := s.bannerRepo.GetBannerList()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *bannerService) AddBanner(params models.Banner) ([]models.Banner, error) {

	now := time.Now()

	bannerAdd := models.Banner{
		Title:          params.Title,
		Subtitle:       params.Subtitle,
		ImageUrl:       params.ImageUrl,
		BannerPosition: params.BannerPosition,
		IsActive:       1,
		CreatedDate:    now,
	}

	data, err := s.bannerRepo.AddBanner(bannerAdd)

	return data, err
}

func (s *bannerService) UpdateBanner(params models.Banner) ([]models.Banner, error) {

	BannerUpdate := models.Banner{
		Id:             params.Id,
		Title:          params.Title,
		Subtitle:       params.Subtitle,
		BannerPosition: params.BannerPosition,
		ImageUrl:       params.ImageUrl,
	}

	data, err := s.bannerRepo.UpdateBanner(BannerUpdate)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *bannerService) UpdateToinactive(paramID int) error {

	err := s.bannerRepo.UpdateToinactive(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (s *bannerService) UpdateToActive(paramID int) error {

	err := s.bannerRepo.UpdateToActive(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (s *bannerService) DeleteBanner(paramID int) error {
	err := s.bannerRepo.DeleteBanner(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
