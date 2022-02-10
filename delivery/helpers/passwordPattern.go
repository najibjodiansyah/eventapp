package helpers

import (
	"errors"
	"regexp"
	"strings"
)

func CheckPasswordPattern(password string) error {
	// Detect blank space in password
	if strings.ContainsAny(password, " ") {
		return errors.New("password contain blank space")
	}

	// Detect insufficient password length
	if len(password) < 6 {
		return errors.New("password must be minimum 6 characters long")
	}

	// Detect absence of lowercase letter
	if re := regexp.MustCompile("[a-z]"); !re.MatchString(password) {
		return errors.New("password must contain lowercase")
	}

	// Detect absence of UPPERCASE letter
	if re := regexp.MustCompile("[A-Z]"); !re.MatchString(password) {
		return errors.New("password must contain UPPERCASE")
	}

	// Detect absence of decimal number
	if re := regexp.MustCompile("[0-9]"); !re.MatchString(password) {
		return errors.New("password must contain decimal number")
	}

	// Detect absence of symbol
	if re := regexp.MustCompile("[~!@#$%^&*]"); !re.MatchString(password) {
		return errors.New("password must contain symbols ~!@#$%^&*")
	}

	return nil
}
