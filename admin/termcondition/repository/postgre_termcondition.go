package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/termcondition"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreSplashScreenRepository struct
type postgreTermConditionRepository struct {
	DbConn *gorm.DB
}

// NewTermConditionRepository DB
func NewTermConditionRepository(DbConn *gorm.DB) termcondition.Repository {
	return &postgreTermConditionRepository{DbConn}
}

func (repo *postgreTermConditionRepository) GetListTermConditionsIsActived() ([]models.TermCondition, error) {
	var termConditions []models.TermCondition
	err := repo.DbConn.Raw(`select * from term_condition where is_actived = 1 and deleted_at is null order by created_at desc`).Scan(&termConditions).Error

	if err != nil {
		return nil, err
	}
	return termConditions, nil
}

func (repo *postgreTermConditionRepository) GetListTermConditions() ([]models.TermCondition, error) {
	var termConditions []models.TermCondition
	err := repo.DbConn.Raw(`select * from term_condition where deleted_at is null order by created_at desc`).Scan(&termConditions).Error

	if err != nil {
		return nil, err
	}
	return termConditions, nil
}

func (repo *postgreTermConditionRepository) Update(term models.TermCondition) ([]models.TermCondition, error) {

	var termCondition []models.TermCondition

	err := repo.DbConn.Table("term_condition").Where("id = ? ", term.ID).Update(&term).Scan(&termCondition).Error
	if err != nil {
		return termCondition, err
	}

	return termCondition, nil

}
