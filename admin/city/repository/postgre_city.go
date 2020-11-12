package repository

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/city"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreCityRepository struct
type postgreCityRepository struct {
	DbConn *gorm.DB
}

// NewCityRepository DB
func NewCityRepository(DbConn *gorm.DB) city.Repository {
	return &postgreCityRepository{DbConn}
}

func (repo *postgreCityRepository) GetCityListAll() ([]models.City, error) {
	var city []models.City

	query := fmt.Sprintf("select * from cities where deleted_at is null order by name asc")

	err := repo.DbConn.Raw(query).Scan(&city).Error

	if err != nil {
		return nil, err
	}
	return city, nil
}

func (repo *postgreCityRepository) AddCity(city models.City) ([]models.City, error) {
	var citySc []models.City

	err := repo.DbConn.Table("cities").Create(&city).Scan(&citySc).Error
	if err != nil {
		return citySc, err
	}

	return citySc, nil
}

func (repo *postgreCityRepository) UpdateCity(city models.City) ([]models.City, error) {

	var citySc []models.City

	err := repo.DbConn.Table("cities").Where("id = ? ", city.ID).Update(&city).Scan(&citySc).Error
	if err != nil {
		return citySc, err
	}
	return citySc, nil
}

func (repo *postgreCityRepository) GetCityByID(paramID int) ([]models.City, error) {

	var City []models.City

	err := repo.DbConn.Table("cities").Where("id = ?", paramID).Scan(&City).Error

	if err != nil {
		return nil, err
	}
	return City, nil

}

func (repo *postgreCityRepository) DeletedCity(paramID int) error {

	var City []models.City
	now := time.Now()

	err := repo.DbConn.Raw(`update cities set deleted_at = ? where id = ?`, now, paramID).Scan(&City).Error

	if err != nil {
		return err
	}
	return err
}

func (repo *postgreCityRepository) GetLastCityID() ([]models.City, error) {
	var city []models.City

	query := fmt.Sprintf("select * from cities order by id desc limit 1")

	err := repo.DbConn.Raw(query).Scan(&city).Error

	if err != nil {
		return nil, err
	}
	return city, nil
}
