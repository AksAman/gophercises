package utils

import (
	"errors"
	"os"
)

func DoesFileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
