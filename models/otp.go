package models

import "time"

type UserOtp struct {
	Id        int
	UserId    int
	OtpCode   string
	IsUsed    int
	CreatedAt time.Time
}

type RequsetValidateOTP struct {
	Email   string `json:"email"`
	OtpCode string `json:"otp_code"`
}
