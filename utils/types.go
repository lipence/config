package utils

import (
	"fmt"
	"math"
)

var ErrInvalidSourceType = fmt.Errorf("invalid source type")
var ErrUnsafeConversion = fmt.Errorf("unsafe conversion")

func ItfToBoolean(src interface{}) (dst bool, err error) {
	switch _src := src.(type) {
	case bool:
		dst = _src
	default:
		return false, ErrInvalidSourceType
	}
	return
}

func ItfToBytes(src interface{}) (dst []byte, err error) {
	if src == nil {
		return nil, nil
	}
	switch _src := src.(type) {
	case []byte:
		dst = _src
	case string:
		dst = []byte(_src)
	case fmt.Stringer:
		dst = []byte(_src.String())
	default:
		return nil, ErrInvalidSourceType
	}
	return
}

func ItfToString(src interface{}) (dst string, err error) {
	if src == nil {
		return "", nil
	}
	switch _src := src.(type) {
	case []byte:
		dst = string(_src)
	case string:
		dst = _src
	case fmt.Stringer:
		dst = _src.String()
	default:
		return "", ErrInvalidSourceType
	}
	return
}

func ItfToStringSlice(src interface{}) (dst []string, err error) {
	if src == nil {
		return nil, nil
	}
	switch _src := src.(type) {
	case []string:
		dst = _src
	case []interface{}:
		dst = make([]string, len(_src))
		for i, v := range _src {
			if vStr, ok := v.(string); !ok {
				return nil, fmt.Errorf("invalid source element %d", i)
			} else {
				dst[i] = vStr
			}
		}
	default:
		return nil, ErrInvalidSourceType
	}
	return
}

func ItfToInt64(src interface{}) (dst int64, err error) {
	switch _src := src.(type) {
	case int:
		return int64(_src), nil
	case int8:
		return int64(_src), nil
	case int16:
		return int64(_src), nil
	case int32:
		return int64(_src), nil
	case int64:
		return _src, nil
	case uint:
		return UInt64ToInt64(uint64(_src))
	case uint8:
		return int64(_src), nil
	case uint16:
		return int64(_src), nil
	case uint32:
		return int64(_src), nil
	case uint64:
		return UInt64ToInt64(_src)
	case float32:
		return Float64ToInt64(float64(_src))
	case float64:
		return Float64ToInt64(_src)
	default:
		return 0, ErrInvalidSourceType
	}
}

func ItfToUInt64(src interface{}) (dst uint64, err error) {
	switch _src := src.(type) {
	case int:
		if _src < 0 {
			return 0, ErrUnsafeConversion
		}
		return uint64(_src), nil
	case int8:
		if _src < 0 {
			return 0, ErrUnsafeConversion
		}
		return uint64(_src), nil
	case int16:
		if _src < 0 {
			return 0, ErrUnsafeConversion
		}
		return uint64(_src), nil
	case int32:
		if _src < 0 {
			return 0, ErrUnsafeConversion
		}
		return uint64(_src), nil
	case int64:
		if _src < 0 {
			return 0, ErrUnsafeConversion
		}
		return uint64(_src), nil
	case uint:
		return uint64(_src), nil
	case uint8:
		return uint64(_src), nil
	case uint16:
		return uint64(_src), nil
	case uint32:
		return uint64(_src), nil
	case uint64:
		return _src, nil
	case float32:
		return Float64ToUInt64(float64(_src))
	case float64:
		return Float64ToUInt64(_src)
	default:
		return 0, ErrInvalidSourceType
	}
}

func ItfToFloat64(src interface{}) (dst float64, err error) {
	switch _src := src.(type) {
	case int:
		return float64(_src), nil
	case int8:
		return float64(_src), nil
	case int16:
		return float64(_src), nil
	case int32:
		return float64(_src), nil
	case int64:
		return float64(_src), nil
	case uint:
		return float64(_src), nil
	case uint8:
		return float64(_src), nil
	case uint16:
		return float64(_src), nil
	case uint32:
		return float64(_src), nil
	case uint64:
		return float64(_src), nil
	case float32:
		return float64(_src), nil
	case float64:
		return _src, nil
	default:
		return 0, ErrInvalidSourceType
	}
}

func UInt64ToInt64(src uint64) (dst int64, err error) {
	if src > uint64(math.MaxInt64) {
		return 0, ErrUnsafeConversion
	}
	return int64(src), nil
}

func Float64ToInt64(src float64) (dst int64, err error) {
	if src > float64(math.MaxInt64) {
		return 0, ErrUnsafeConversion
	}
	if float64(int64(src)) != src {
		return 0, ErrUnsafeConversion
	}
	return int64(src), nil
}

func Float64ToUInt64(src float64) (dst uint64, err error) {
	if src < 0 || src > float64(math.MaxUint64) {
		return 0, ErrUnsafeConversion
	}
	if float64(int64(src)) != src {
		return 0, ErrUnsafeConversion
	}
	return uint64(src), nil
}
