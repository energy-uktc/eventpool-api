package models

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("changepasswordrequest", changePasswordValidate)
	}
}

type CreateUserRequest struct {
	UserName string `json:"userName" binding:"required"`
	AuthUserRequest
}

type AuthUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUserRequest struct {
	VerificationCode  string `json:"verificationCode" binding:"required"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type VerifyUserResponse struct {
	Email string
	Token *GeneratedTokenResponse
}

type RefreshTokenRequest struct {
	Token        string `json:"token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type GeneratedToken struct {
	Token          string `json:"token"`
	TokenType      string `json:"token_type"`
	TokenId        string `json:"jti"`
	ExpirationTime int64  `json:"expiration_time"`
	Scope          string `json:"scope"`
}

type GeneratedTokenResponse struct {
	GeneratedToken
	RefreshToken string `json:"refresh_token"`
}

type UserContextInfo struct {
	UserId string
	Scopes []string
}

type ChangePasswordRequest struct {
	Action           string `json:"action" binding:"required,changepasswordrequest"`
	Email            string `json:"email"`
	VerificationCode string `json:"verificationCode"`
	OldPassword      string `json:"oldPassword"`
	NewPassword      string `json:"newPassword"`
}

var changePasswordValidate validator.Func = func(fl validator.FieldLevel) bool {
	var request ChangePasswordRequest
	request = fl.Parent().Interface().(ChangePasswordRequest)
	action := fl.Field().Interface().(string)
	switch action {
	case "change":
		if request.OldPassword == "" || request.NewPassword == "" || request.Email == "" {
			return false
		}
		return true
	case "reset":
		if request.VerificationCode == "" || request.NewPassword == "" {
			return false
		}
		return true
	case "sendResetEmail":
		if request.Email == "" {
			return false
		}
		return true
	default:
		return false
	}
}
