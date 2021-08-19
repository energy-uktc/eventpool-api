package user_repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/energy-uktc/eventpool-api/database"
	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func Create(user *entities.User) error {
	response := database.DbConn.Create(user)
	if response.Error != nil {
		log.Println(response.Error)
		return fmt.Errorf("Something went wrong during registration. PLease try again!")
	}
	return nil
}

func FindById(id string) (*entities.User, error) {
	var existingUser *entities.User
	response := database.DbConn.First(&existingUser, uuid.FromStringOrNil(id))
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User Not Found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return existingUser, nil
}

func FindByEmail(userEmail string) (*entities.User, error) {
	var existingUser *entities.User
	response := database.DbConn.First(&existingUser, "email = ?", userEmail)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return existingUser, nil
}

func FindByVerificationCode(verificationCode string) (*entities.User, error) {
	var userToken *entities.UserToken
	response := database.DbConn.Joins("User").First(&userToken, "type = ? AND token = ?", entities.VerificationCode, verificationCode)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Verification code not found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return &userToken.User, nil
}

func FindRefreshToken(token string) (*entities.UserToken, error) {
	var userToken *entities.UserToken
	response := database.DbConn.Joins("User").First(&userToken, "type = ? AND token = ?", entities.RefreshToken, token)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Refresh token not found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return userToken, nil
}

func DeleteVerificationCodes(user *entities.User) error {
	response := database.DbConn.Delete(entities.UserToken{}, "user_id = ? AND type = ?", user.ID, entities.VerificationCode)
	if response.Error != nil {
		log.Println(response.Error)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func DeleteRefreshTokens(user *entities.User) error {
	response := database.DbConn.Delete(entities.UserToken{}, "user_id = ? AND type = ?", user.ID, entities.RefreshToken)
	if response.Error != nil {
		log.Println(response.Error)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func CreateVerificationToken(token string, user *entities.User) (*entities.UserToken, error) {
	userToken, err := createUserToken(entities.VerificationCode, token, user)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Something went wrong during verification code creation. PLease try again!")
	}
	return userToken, nil
}

func CreateRefreshToken(token string, user *entities.User) (*entities.UserToken, error) {
	userToken, err := createUserToken(entities.RefreshToken, token, user)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Something went wrong. PLease try again!")
	}
	return userToken, nil
}

func createUserToken(tokenType entities.TokenType, token string, user *entities.User) (*entities.UserToken, error) {
	userToken := &entities.UserToken{
		UserID: user.ID,
		Type:   tokenType,
		Token:  token,
	}
	if err := database.DbConn.Create(&userToken).Error; err != nil {
		return nil, err
	}
	return userToken, nil
}

func SetUserVerified(user *entities.User) error {
	return database.DbConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Updates(map[string]interface{}{"Verified": true}).Error; err != nil {
			log.Println(err)
			return fmt.Errorf("Something went wrong during user verification. PLease try again!")
		}
		return nil
	})
}
