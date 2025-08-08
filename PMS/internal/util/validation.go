package util

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func ValidateEmail(email string) error {
	Email:=strings.ToLower(email)
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	if !re.MatchString(Email){
         return errors.New("invalid email address")
	}
	return nil
}

func ValidateMobileNumber(phone string) error {

	re := regexp.MustCompile(`^[6-9]\d{9}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number: must be 10 digits and start with 6-9")
	}
	return nil
}
func ValidatePassword(password string)error{
	if len(password) < 12{
		return errors.New("password must be at least 12 charactercontain at least one uppercase letter, one lowercase letter, one number, and one special symbol")

	}

	var hasUpper, hasLowe, hasDigit , hasSpecial bool

	for _,char:=range password{
		switch{
		case unicode.IsUpper(char):
			hasUpper=true
		case unicode.IsLower(char):
			hasLowe=true
		case unicode.IsDigit(char):
			hasDigit=true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial=true

		}
	}
	if !(hasDigit && hasLowe && hasUpper && hasSpecial){
		return errors.New("password must be at least 12 charactercontain at least one uppercase letter, one lowercase letter, one number, and one special symbol")
	}
	return nil
}

