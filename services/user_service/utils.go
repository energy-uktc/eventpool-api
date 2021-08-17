package user_service

import (
	"log"

	"github.com/energy-uktc/eventpool-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func getVerificationCode(max int) string {
	return utils.GenerateString(max, utils.ALPHA_NUMERIC)
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comaprePassword(provided string, compareWith string) error {
	return bcrypt.CompareHashAndPassword([]byte(compareWith), []byte(provided))
}
