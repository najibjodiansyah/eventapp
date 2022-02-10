package helpers

import (
	"errors"
	"time"
)

func CheckDatetimePattern(dt string) error {
	layout := "2006-01-02 15:04:05.000"

	if _, err := time.Parse(layout, dt); err != nil {
		return errors.New("incorrect datetime format")
	}

	return nil
}
