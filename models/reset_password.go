package models

type RequestResetPassword struct {
	NewPassword string `json:"new_password"`
	Email       string `json:"email"`
	OtpCode     string `json:"otp_code"`
}
