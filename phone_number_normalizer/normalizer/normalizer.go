package normalizer

import (
	"strings"
	"unicode"
)

func NormalizePhoneNumber(phoneNumber string) string {
	var normalizedNumberBuilder strings.Builder

	for _, char := range phoneNumber {
		if unicode.IsDigit(char) {
			normalizedNumberBuilder.WriteRune(char)
		}
	}

	return normalizedNumberBuilder.String()
}
