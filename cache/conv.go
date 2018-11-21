package cache

import (
	"fmt"
	"strconv"
)

// GetString convert interface to string.
func GetString(v interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}

	switch result := v.(type) {
	case string:
		return result, nil
	case []byte:
		return string(result), nil
	default:
		if v != nil {
			return fmt.Sprint(result), nil
		}
	}
	return "", nil
}

// GetInt convert interface to int.
func GetInt(v interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	switch result := v.(type) {
	case int:
		return result, nil
	case int32:
		return int(result), nil
	case int64:
		return int(result), nil
	default:
		if d, _ := GetString(v, nil); d != "" {
			value, _ := strconv.Atoi(d)
			return value, nil
		}
	}
	return 0, nil
}

// GetInt64 convert interface to int64.
func GetInt64(v interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	switch result := v.(type) {
	case int:
		return int64(result), nil
	case int32:
		return int64(result), nil
	case int64:
		return result, nil
	default:

		if d, _ := GetString(v, nil); d != "" {
			value, _ := strconv.ParseInt(d, 10, 64)
			return value, nil
		}
	}
	return 0, nil
}

// GetFloat64 convert interface to float64.
func GetFloat64(v interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}

	switch result := v.(type) {
	case float64:
		return result, nil
	default:
		if d, _ := GetString(v, nil); d != "" {
			value, _ := strconv.ParseFloat(d, 64)
			return value, nil
		}
	}
	return 0, nil
}

// GetBool convert interface to bool.
func GetBool(v interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}

	switch result := v.(type) {
	case bool:
		return result, nil
	default:
		if d, _ := GetString(v, nil); d != "" {
			value, _ := strconv.ParseBool(d)
			return value, nil
		}
	}
	return false, nil
}
