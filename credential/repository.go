package credential

import "github.com/rudyraharjo/emurojaah/models"

type Repository interface {

	//Get list
	GetUserIdByAlias(alias string) (int, error)
	GetUserBasicInfoById(userId int) (*models.UserBasicInfo, error)
	GetLatestSentOTP(userId int) (*models.UserOtp, error)

	//Insert
	InsertUserOTP(userOTP models.UserOtp) error

	//Update
	UpdateUserOtpAsUsed(id int) error
	UpdateUserPasswordById(userId int, newPassword string) error
}
