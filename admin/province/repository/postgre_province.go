package repository

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/province"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreProvinceRepository struct
type postgreProvinceRepository struct {
	DbConn *gorm.DB
}

// NewProvinceRepository DB
func NewProvinceRepository(DbConn *gorm.DB) province.Repository {
	return &postgreProvinceRepository{DbConn}
}

func (repo *postgreProvinceRepository) GetProvinceListAll() ([]models.Province, error) {
	var province []models.Province

	query := fmt.Sprintf("select * from provinces where deleted_at is null order by name asc")

	err := repo.DbConn.Raw(query).Scan(&province).Error

	if err != nil {
		return nil, err
	}
	return province, nil
}

func (repo *postgreProvinceRepository) AddProvince(province models.Province) ([]models.Province, error) {
	var provinceSc []models.Province

	err := repo.DbConn.Table("provinces").Create(&province).Scan(&provinceSc).Error
	if err != nil {
		return provinceSc, err
	}

	return provinceSc, nil
}

func (repo *postgreProvinceRepository) UpdateProvince(province models.Province) ([]models.Province, error) {

	var provinceSc []models.Province

	err := repo.DbConn.Table("provinces").Where("id = ? ", province.ID).Update(&province).Scan(&provinceSc).Error
	if err != nil {
		return provinceSc, err
	}
	return provinceSc, nil
}

func (repo *postgreProvinceRepository) GetProvinceByID(paramID int) ([]models.Province, error) {

	var Province []models.Province

	err := repo.DbConn.Table("provinces").Where("id = ?", paramID).Scan(&Province).Error

	if err != nil {
		return nil, err
	}
	return Province, nil

}

func (repo *postgreProvinceRepository) DeletedProvince(paramID int) error {

	var Province []models.Province
	now := time.Now()

	err := repo.DbConn.Raw(`update provinces set deleted_at = ? where id = ?`, now, paramID).Scan(&Province).Error

	if err != nil {
		return err
	}
	return err
}

func (repo *postgreProvinceRepository) GetLastProvinceID() ([]models.Province, error) {
	var province []models.Province

	query := fmt.Sprintf("select * from provinces order by id desc limit 1")

	err := repo.DbConn.Raw(query).Scan(&province).Error

	if err != nil {
		return nil, err
	}
	return province, nil
}
