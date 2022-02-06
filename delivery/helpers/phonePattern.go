package helpers

import (
	"errors"
	"regexp"
)

func CheckPhonePattern(phone string) error {
	re := regexp.MustCompile("[^+0-9]")

	if re.MatchString(phone) {
		return errors.New("phone number contain non numeric")
	}

	return nil
}
