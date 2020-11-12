package http

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/credential"
	"github.com/rudyraharjo/emurojaah/models"
)

type HttpCredentialHandler struct {
	CredentialService credential.Service
}

func NewCredentialHttpHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, credentialService credential.Service) {
	handler := HttpCredentialHandler{credentialService}

	r.POST("/api/credential/forgot-password", handler.HandlerForgotPassword)
	r.POST("/api/credential/otp/validate", handler.HandlerValidateOtp)
	r.POST("/api/credential/reset-password", handler.HandlerResetPassword)
}

func (h *HttpCredentialHandler) HandlerForgotPassword(ctx *gin.Context) {
	var params models.RequestForgotPassword
	_ = ctx.Bind(&params)

	if params.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	userId := h.CredentialService.GetUserIdByAlias(params.Email)
	if userId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	userInfo, err1 := h.CredentialService.GetUserBasicInfoById(userId)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error get basic user info",
		})
		return
	}

	otpCode := h.CredentialService.GenerateOTPCode(6)
	err2 := h.CredentialService.InsertUserOTP(userId, otpCode)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error insert otp code",
		})
		return
	}

	err3 := h.CredentialService.SendEmailOTPForgotPassword(ctx, userInfo.FullName, params.Email, otpCode)
	if err3 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error send email",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "OTP code send",
	})

}

func (h *HttpCredentialHandler) HandlerValidateOtp(ctx *gin.Context) {
	var params models.RequsetValidateOTP
	_ = ctx.Bind(&params)

	if params.Email == "" || params.OtpCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	userId := h.CredentialService.GetUserIdByAlias(params.Email)
	latestOtp := h.CredentialService.GetLatestSentOtp(userId)

	if latestOtp == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusNotAcceptable,
			"message": "user has no otp",
		})
		return
	}

	if params.OtpCode == latestOtp.OtpCode {

		if latestOtp.IsUsed == 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusLocked,
				"message": "OTP already used",
			})
			return
		}

		currentTime := time.Now().UnixNano() / 1000000
		expiredTime := latestOtp.CreatedAt.Add(5*time.Minute).UnixNano() / 1000000

		if currentTime > expiredTime {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusGone,
				"message": "OTP expired",
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "OTP match",
			})
			return
		}

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusConflict,
			"message": "OTP Code didn't match",
		})
		return
	}
}

func (h *HttpCredentialHandler) HandlerResetPassword(ctx *gin.Context) {
	var params models.RequestResetPassword
	_ = ctx.Bind(&params)

	if params.NewPassword == "" || params.Email == "" || params.OtpCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	userId := h.CredentialService.GetUserIdByAlias(params.Email)
	latestOtp := h.CredentialService.GetLatestSentOtp(userId)

	if latestOtp == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusNotAcceptable,
			"message": "user has no otp",
		})
		return
	}

	if params.OtpCode == latestOtp.OtpCode {

		if latestOtp.IsUsed == 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusLocked,
				"message": "OTP already used",
			})
			return
		}

		currentTime := time.Now().UnixNano() / 1000000
		expiredTime := latestOtp.CreatedAt.Add(5*time.Minute).UnixNano() / 1000000

		if currentTime > expiredTime {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusGone,
				"message": "OTP expired",
			})
			return
		} else {

			hashedPassword, _ := h.CredentialService.HashPassword(params.NewPassword)
			errUpdatePass := h.CredentialService.UpdateUserPasswordById(userId, hashedPassword)

			if errUpdatePass != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "system error on update password",
				})
				return
			}

			// set otp as used
			h.CredentialService.UpdateUserOtpAsUsed(latestOtp.Id)

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Password changed!",
			})
			return
		}

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusConflict,
			"message": "OTP Code didn't match",
		})
		return
	}

}
