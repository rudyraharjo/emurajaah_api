package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/misc"
	"github.com/rudyraharjo/emurojaah/models"
)

type postgreMiscRepository struct {
	DbCon *gorm.DB
}

func NewMiscRepository(db *gorm.DB) misc.Repository {
	return &postgreMiscRepository{db}
}

func (repo *postgreMiscRepository) GetProvinceList() ([]models.Province, error) {
	var province []models.Province
	err := repo.DbCon.Table("dukcapil_propinsi").Find(&province).Error
	if err != nil {
		return nil, err
	}
	return province, nil
}
