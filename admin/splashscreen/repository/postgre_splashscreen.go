package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/splashscreen"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreSplashScreenRepository struct
type postgreSplashScreenRepository struct {
	DbConn *gorm.DB
}

// NewSplashScreenRepository DB
func NewSplashScreenRepository(DbConn *gorm.DB) splashscreen.Repository {
	return &postgreSplashScreenRepository{DbConn}
}

func (repo *postgreSplashScreenRepository) GetSplashScreenList() ([]models.SplashScreen, error) {
	var splashScreen []models.SplashScreen

	query := fmt.Sprintf("select * from splash_screen where is_active = 1 order by position asc")

	err := repo.DbConn.Raw(query).Find(&splashScreen).Error

	if err != nil {
		return nil, err
	}
	return splashScreen, nil
}

func (repo *postgreSplashScreenRepository) GetSplashScreenListAll() ([]models.SplashScreen, error) {
	var splashScreen []models.SplashScreen

	query := fmt.Sprintf("select * from splash_screen order by position asc")

	err := repo.DbConn.Raw(query).Find(&splashScreen).Error

	if err != nil {
		return nil, err
	}
	return splashScreen, nil
}

func (repo *postgreSplashScreenRepository) GetSplashScreenListActive() ([]models.SplashScreen, error) {
	var splashScreen []models.SplashScreen

	query := fmt.Sprintf("select * from splash_screen where is_active = 1 order by position asc")

	err := repo.DbConn.Raw(query).Find(&splashScreen).Error

	if err != nil {
		return nil, err
	}
	return splashScreen, nil
}

func (repo *postgreSplashScreenRepository) AddSplashScreen(splashscreen models.SplashScreen) ([]models.SplashScreen, error) {

	var SplashSc []models.SplashScreen

	err := repo.DbConn.Table("splash_screen").Create(&splashscreen).Scan(&SplashSc).Error
	if err != nil {
		return SplashSc, err
	}

	return SplashSc, nil

}

func (repo *postgreSplashScreenRepository) UpdateSplashScreen(splash models.SplashScreen) ([]models.SplashScreen, error) {

	var SplashSc []models.SplashScreen

	err := repo.DbConn.Table("splash_screen").Where("id = ? ", splash.ID).Update(&splash).Scan(&SplashSc).Error
	if err != nil {
		fmt.Print(err)
		return SplashSc, err
	}
	return SplashSc, nil
}

func (repo *postgreSplashScreenRepository) UpdateToinactive(paramID int) error {

	err := repo.DbConn.Table("splash_screen").Where("id = ? ", paramID).Update("is_active", 0).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *postgreSplashScreenRepository) UpdateToActive(paramID int) error {

	err := repo.DbConn.Table("splash_screen").Where("id = ? ", paramID).Update("is_active", 1).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *postgreSplashScreenRepository) DeleteSplashScreen(ID int) error {

	query := fmt.Sprintf(`DELETE FROM splash_screen WHERE id = $1;`)

	err := repo.DbConn.Exec(query, ID).Error
	if err != nil {
		return err
	}

	return nil
}
