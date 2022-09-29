package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetQueryParam[T comparable](c *fiber.Ctx, key string, defaultValue T) T {

	value := c.Query(key)

	if value == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case string:
		return any(value).(T)
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T)
		}
		return defaultValue
	case bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return any(boolValue).(T)
		}
		return defaultValue
	default:
		return defaultValue
	}

}
