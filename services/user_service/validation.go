package user_service

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"
)

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("Password not valid. Must contain minimum eight characters.")
	}
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
		case unicode.IsUpper(ch):
			uppercasePresent = true
		case unicode.IsLower(ch):
			lowercasePresent = true
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("at least one numeric character required")
	}
	if len(errorString) != 0 {
		errorString = "Password not valid: " + errorString
		return fmt.Errorf(errorString)
	}
	return nil
}

func validateEmail(email string) error {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return fmt.Errorf("Email not valid.Must be between 4 and 254 characters")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("Email not valid")
	}
	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return fmt.Errorf("Email not valid")
	}
	return nil
}
