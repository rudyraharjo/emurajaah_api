package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/point"
	"github.com/rudyraharjo/emurojaah/models"
)

type postgrePointRepository struct {
	DbConn *gorm.DB
}

// NewPointRepository DB
func NewPointRepository(DbConn *gorm.DB) point.Repository {
	return &postgrePointRepository{DbConn}
}

func (repo *postgrePointRepository) GetPointList() ([]models.ResponsePoint, error) {

	var points []models.ResponsePoint

	query := fmt.Sprintf("select * from point where is_active = 1 order by id asc")

	err := repo.DbConn.Raw(query).Find(&points).Error

	if err != nil {
		return nil, err
	}
	return points, nil
}

func (repo *postgrePointRepository) AddPoint(point models.ResponsePoint) ([]models.ResponsePoint, error) {
	var insertPoint []models.ResponsePoint

	err := repo.DbConn.Table("point").Create(&point).Scan(&insertPoint).Error
	if err != nil {
		return nil, err
	}

	return insertPoint, nil
}

func (repo *postgrePointRepository) UpdatePoint(point models.ResponsePoint) ([]models.ResponsePoint, error) {

	var updatePoint []models.ResponsePoint

	err := repo.DbConn.Table("point").Where("id = ?", point.ID).Update(&point).Scan(&updatePoint).Error

	if err != nil {
		return nil, err
	}
	return updatePoint, nil

}

func (repo *postgrePointRepository) UpdateToinactive(paramID int) error {

	err := repo.DbConn.Table("point").Where("id = ? ", paramID).Update("is_active", 0).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *postgrePointRepository) UpdateToActive(paramID int) error {

	err := repo.DbConn.Table("point").Where("id = ? ", paramID).Update("is_active", 1).Error

	if err != nil {
		return err
	}

	return nil
}
