package iz2go

import (
	"reflect"
	"strconv"
)

func ParseOrDefault[T any](value any, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	if v, ok := value.(T); ok {
		return v
	}

	switch any(defaultValue).(type) {
	case string:
		return any(ParseString(value)).(T)
	case int, int8, int16, int32, int64:
		return any(ParseInt(value)).(T)
	case float32, float64:
		return any(ParseFloat(value)).(T)
	case bool:
		return any(ParseBool(value)).(T)
	}

	return defaultValue
}

func ParseString(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	}
	return ""
}

func ParseInt(value any) int {
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return i
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(v).Int())
	case float32, float64:
		return int(reflect.ValueOf(v).Float())
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func ParseFloat(value any) float64 {
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return f
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int())
	case float32, float64:
		return reflect.ValueOf(v).Float()
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

func ParseBool(value any) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return b
	case int, int8, int16, int32, int64:
		return v != 0
	case float32, float64:
		return v != 0
	}
	return false
}
