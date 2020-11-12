package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/banner"
	"github.com/rudyraharjo/emurojaah/models"
)

type postgreBannerRepository struct {
	DbConn *gorm.DB
}

// NewBannerRepository DB
func NewBannerRepository(DbConn *gorm.DB) banner.Repository {
	return &postgreBannerRepository{DbConn}
}

func (repo *postgreBannerRepository) GetBannerList() ([]models.Banner, error) {

	var banners []models.Banner

	query := fmt.Sprintf("select * from banner order by id asc")

	err := repo.DbConn.Raw(query).Find(&banners).Error

	if err != nil {
		return nil, err
	}
	return banners, nil
}

func (repo *postgreBannerRepository) AddBanner(banner models.Banner) ([]models.Banner, error) {
	var insertBanner []models.Banner

	err := repo.DbConn.Table("banner").Create(&banner).Scan(&insertBanner).Error
	if err != nil {
		return nil, err
	}

	return insertBanner, nil
}

func (repo *postgreBannerRepository) UpdateBanner(banner models.Banner) ([]models.Banner, error) {

	var updateBanner []models.Banner

	err := repo.DbConn.Table("banner").Where("id = ?", banner.Id).Update(&banner).Scan(&updateBanner).Error

	if err != nil {
		return nil, err
	}
	return updateBanner, nil

}

func (repo *postgreBannerRepository) UpdateToinactive(paramID int) error {

	err := repo.DbConn.Table("banner").Where("id = ? ", paramID).Update("is_active", 0).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *postgreBannerRepository) UpdateToActive(paramID int) error {

	err := repo.DbConn.Table("banner").Where("id = ? ", paramID).Update("is_active", 1).Error

	if err != nil {
		return err
	}

	return nil
}

//DeleteBanner
func (repo *postgreBannerRepository) DeleteBanner(paramID int) error {

	query := fmt.Sprintf(`DELETE FROM banner WHERE id = $1;`)

	err := repo.DbConn.Exec(query, paramID).Error
	if err != nil {
		return err
	}

	return nil
}

//GetBannerListActive
func (repo *postgreBannerRepository) GetBannerListActive() ([]models.Banner, error) {
	var banners []models.Banner

	query := fmt.Sprintf("select * from banner where is_active = 1 order by id asc")

	err := repo.DbConn.Raw(query).Find(&banners).Error

	if err != nil {
		return nil, err
	}
	return banners, nil
}
