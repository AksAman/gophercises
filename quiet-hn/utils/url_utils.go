package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetQueryParam[T comparable](c echo.Context, key string, defaultValue T) T {

	value := c.QueryParam(key)

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
