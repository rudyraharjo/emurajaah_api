package credential

import (
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/models"
)

type Service interface {

	//Get List
	GetUserIdByAlias(alias string) int
	GetUserBasicInfoById(userId int) (*models.UserBasicInfo, error)
	GetLatestSentOtp(userId int) *models.UserOtp

	//Insert
	InsertUserOTP(userId int, otpCode string) error

	//Update
	UpdateUserOtpAsUsed(id int)
	UpdateUserPasswordById(userId int, newPassword string) error

	// Other
	HashPassword(password string) (string, error)
	GenerateOTPCode(max int) string
	SendEmailOTPForgotPassword(c *gin.Context, userName string, userEmail string, otpCode string) error
}
