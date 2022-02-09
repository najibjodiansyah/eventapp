package helpers

import (
	"errors"
	"strings"
)

func CheckStringInput(s string) error {
	if strings.ContainsAny(s, ";") {
		return errors.New("input cannot contain forbidden character")
	}

	return nil
}
