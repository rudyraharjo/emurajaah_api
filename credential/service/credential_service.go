package service

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/credential"
	"github.com/rudyraharjo/emurojaah/models"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

type credentialService struct {
	credentialRepo credential.Repository
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func NewCredentailService(cr credential.Repository) credential.Service {
	return &credentialService{
		cr,
	}
}

//HashPassword - Encrypt password
func (s *credentialService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *credentialService) GetUserIdByAlias(alias string) int {
	id, err := s.credentialRepo.GetUserIdByAlias(alias)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return id
}

func (s *credentialService) GetUserBasicInfoById(userId int) (*models.UserBasicInfo, error) {
	usr, err := s.credentialRepo.GetUserBasicInfoById(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return usr, nil
}

func (s *credentialService) GenerateOTPCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func (s *credentialService) SendEmailOTPForgotPassword(c *gin.Context, userName string, userEmail string, otpCode string) error {

	const confirmMessage = `
	Hi %s,<br>
	Berikut kode OTP untuk melakukan reset password: <br><br>

	<b>%s</b><br><br>
	
	Jangan berikan kode OTP kepada siapapun, termasuk pihak yang mengatasnamakan eMurojaah.`

	from := mail.NewEmail("eMurojaah", "no-reply@emurojaah.com")
	subject := "Kode OTP Reset Password"
	to := mail.NewEmail(userName, userEmail)
	message := mail.NewSingleEmail(from, subject, to, fmt.Sprintf(confirmMessage, userName, otpCode), fmt.Sprintf(confirmMessage, userName, otpCode))
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *credentialService) InsertUserOTP(userId int, otpCode string) error {
	userOTP := models.UserOtp{
		UserId:    userId,
		OtpCode:   otpCode,
		IsUsed:    0,
		CreatedAt: time.Now().UTC(),
	}

	err := s.credentialRepo.InsertUserOTP(userOTP)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *credentialService) GetLatestSentOtp(userId int) *models.UserOtp {
	userOtp, err := s.credentialRepo.GetLatestSentOTP(userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return userOtp
}

func (s *credentialService) UpdateUserOtpAsUsed(id int) {
	err := s.credentialRepo.UpdateUserOtpAsUsed(id)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *credentialService) UpdateUserPasswordById(userId int, newPassword string) error {
	err := s.credentialRepo.UpdateUserPasswordById(userId, newPassword)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
