package models

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
