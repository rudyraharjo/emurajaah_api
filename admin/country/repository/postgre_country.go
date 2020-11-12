package repository

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/country"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreCountryRepository struct
type postgreCountryRepository struct {
	DbConn *gorm.DB
}

// NewCountryRepository DB
func NewCountryRepository(DbConn *gorm.DB) country.Repository {
	return &postgreCountryRepository{DbConn}
}

func (repo *postgreCountryRepository) GetCountryListAll() ([]models.Country, error) {
	var quote []models.Country

	query := fmt.Sprintf("select * from country where deleted_at is null order by created_at desc")

	err := repo.DbConn.Raw(query).Scan(&quote).Error

	if err != nil {
		return nil, err
	}
	return quote, nil
}

func (repo *postgreCountryRepository) AddCountry(country models.Country) ([]models.Country, error) {
	var countrySc []models.Country

	err := repo.DbConn.Table("country").Create(&country).Scan(&countrySc).Error
	if err != nil {
		return countrySc, err
	}

	return countrySc, nil
}

func (repo *postgreCountryRepository) UpdateCountry(country models.Country) ([]models.Country, error) {

	var countrySc []models.Country

	err := repo.DbConn.Table("country").Where("id = ? ", country.ID).Update(&country).Scan(&countrySc).Error
	if err != nil {
		return countrySc, err
	}
	return countrySc, nil

}

func (repo *postgreCountryRepository) GetCountryByID(paramID int) ([]models.Country, error) {

	var Country []models.Country

	err := repo.DbConn.Table("country").Where("id = ?", paramID).Scan(&Country).Error

	if err != nil {
		return nil, err
	}
	return Country, nil

}

func (repo *postgreCountryRepository) DeletedCountry(paramID int) error {

	var Country []models.Country
	now := time.Now()

	err := repo.DbConn.Raw(`update country set deleted_at = ? where id = ?`, now, paramID).Scan(&Country).Error

	if err != nil {
		return err
	}
	return err
}
