package user_service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/repositories/user_repository"
	"github.com/energy-uktc/eventpool-api/services/jwt_service"
	"github.com/energy-uktc/eventpool-api/utils"
)

func RegisterUser(userModel *models.CreateUserRequest) (*entities.User, error) {
	if err := validateEmail(userModel.Email); err != nil {
		return nil, err
	}
	if err := validatePassword(userModel.Password); err != nil {
		return nil, err
	}
	user, err := user_repository.FindByEmail(userModel.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("Email already exists!")
	}
	user = &entities.User{UserName: userModel.UserName, Email: userModel.Email, Password: hashAndSalt(userModel.Password)}
	err = user_repository.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUnverifiedUser(userModel *models.AuthUserRequest) (*entities.User, error) {
	user, _ := user_repository.FindByEmail(userModel.Email)
	if user == nil {
		return nil, fmt.Errorf("User %s not found", userModel.Email)
	}
	err := checkPasswordMatch(userModel.Password, user)
	if err != nil {
		return nil, err
	}
	if user.Verified {
		return nil, fmt.Errorf("User is already verified")
	}
	return user, nil
}

func VerifyUserCode(verifyCodeModel *models.VerifyUserRequest) (*models.VerifyUserResponse, error) {
	user, err := user_repository.FindByVerificationCode(verifyCodeModel.VerificationCode)
	if err != nil {
		return nil, err
	}
	if user.Verified {
		return nil, fmt.Errorf("User is already verified")
	}

	if err := user_repository.SetUserVerified(user); err != nil {
		return nil, err
	}

	if verifyCodeModel.ReturnSecureToken {
		tokenResponse, err := generateToken(user)
		if err != nil {
			return nil, err
		}
		return &models.VerifyUserResponse{
			Email: user.Email,
			Token: tokenResponse,
		}, nil
	}
	return nil, nil
}

func GenerateToken(userRequest *models.AuthUserRequest) (*models.GeneratedTokenResponse, error) {
	user, _ := user_repository.FindByEmail(userRequest.Email)
	if user == nil {
		return nil, fmt.Errorf("User %s not found", userRequest.Email)
	}

	if !user.Verified {
		return nil, fmt.Errorf("User is not verified")
	}

	if err := checkPasswordMatch(userRequest.Password, user); err != nil {
		return nil, err
	}

	return generateToken(user)

}

func RefreshToken(refreshTokenRequest *models.RefreshTokenRequest) (*models.GeneratedTokenResponse, error) {
	token, err := jwt_service.VerifyToken(refreshTokenRequest.Token, true)
	if err != nil {
		return nil, err
	}
	claims := jwt_service.GetClaims(token)
	hashedRefreshToken := sha256.Sum256([]byte(refreshTokenRequest.RefreshToken + claims.Id))
	userToken, _ := user_repository.FindRefreshToken(hex.EncodeToString(hashedRefreshToken[:]))
	if userToken == nil {
		return nil, fmt.Errorf("Refresh token not valid")
	}

	if userToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("Refresh token has expired")
	}

	if claims.CustomerInfo.Id != userToken.UserID.String() {
		return nil, fmt.Errorf("Refresh token not valid")
	}

	return generateToken(&userToken.User)
}

func RevokeRefreshToken(userContext *models.UserContextInfo) error {
	return user_repository.DeleteRefreshTokens(entities.NewUserFromID(userContext.UserId))
}

func checkPasswordMatch(password string, user *entities.User) error {
	err := comaprePassword(password, user.Password)
	if err != nil {
		return fmt.Errorf("Incorrect Password")
	}
	return nil
}

func GenerateVerificationCode(user *entities.User) (string, error) {
	err := user_repository.DeleteVerificationCodes(user)
	if err != nil {
		return "", err
	}
	verificationCode := getVerificationCode(10)
	createdCode, err := user_repository.CreateVerificationToken(verificationCode, user)
	if err != nil {
		return "", err
	}
	return createdCode.Token, nil
}

func generateToken(user *entities.User) (*models.GeneratedTokenResponse, error) {
	generatedToken, err := jwt_service.CreateToken(user.ID.String(), nil)
	if err != nil {
		return nil, err
	}

	refreshToken := utils.GenerateStringRandomLength(128, 256)
	hashedRefreshToken := sha256.Sum256([]byte(refreshToken + generatedToken.TokenId))
	err = user_repository.DeleteRefreshTokens(user)
	if err != nil {
		return nil, err
	}
	if _, err := user_repository.CreateRefreshToken(hex.EncodeToString(hashedRefreshToken[:]), user); err != nil {
		return nil, err
	}

	return &models.GeneratedTokenResponse{
		GeneratedToken: generatedToken,
		RefreshToken:   refreshToken,
	}, nil
}