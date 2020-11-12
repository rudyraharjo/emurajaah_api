package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/credential"
	"github.com/rudyraharjo/emurojaah/models"
)

type postgreCredentialRepository struct {
	DbConn *gorm.DB
}

func NewCredentialRepository(DbConn *gorm.DB) credential.Repository {
	return &postgreCredentialRepository{
		DbConn,
	}
}

func (repo *postgreCredentialRepository) GetUserIdByAlias(alias string) (int, error) {
	var userAlias models.UserAliasFull
	err := repo.DbConn.Table("user_alias").Where("alias = ?", alias).First(&userAlias).Error
	if err != nil {
		return 0, err
	}
	return userAlias.UserId, nil
}

func (repo *postgreCredentialRepository) InsertUserOTP(userOTP models.UserOtp) error {
	err := repo.DbConn.Table("user_otp").Create(&userOTP).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreCredentialRepository) GetUserBasicInfoById(userId int) (*models.UserBasicInfo, error) {
	var userBasic models.UserBasicInfo
	err := repo.DbConn.Table("user").Where("id = ?", userId).Find(&userBasic).Error
	if err != nil {
		return nil, err
	}
	return &userBasic, nil
}

func (repo *postgreCredentialRepository) GetLatestSentOTP(userId int) (*models.UserOtp, error) {
	var userOtp models.UserOtp
	err := repo.DbConn.Table("user_otp").Where("user_id = ?", userId).Order("created_at desc").First(&userOtp).Error

	if err != nil {
		return nil, err
	}
	return &userOtp, nil
}

func (repo *postgreCredentialRepository) UpdateUserOtpAsUsed(id int) error {
	err := repo.DbConn.Table("user_otp").Where("id = ?", id).Update("is_used", 1).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreCredentialRepository) UpdateUserPasswordById(userId int, newPassword string) error {
	err := repo.DbConn.Table("user_alias").Where("user_id = ?", userId).Update("credential", newPassword).Error
	if err != nil {
		return err
	}
	return nil
}
