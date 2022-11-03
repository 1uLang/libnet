package types

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Byte 将值转换成byte
func Byte(value interface{}) byte {
	return Uint8(value)
}

// Int 将值转换成int
func Int(value interface{}) int {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		return x
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case uint:
		return int(x)
	case uint8:
		return int(x)
	case uint16:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)
	case float32:
		return int(x)
	case float64:
		return int(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 32)
		if err == nil {
			return int(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 32)
			if err == nil {
				return int(floatResult)
			}
		}
	case []byte:
		return Int(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Int(v.Float())
			} else if unsigned {
				return Int(v.Uint())
			} else {
				return Int(v.Int())
			}
		}
	}
	return 0
}

// Int8 将值转换成int8
func Int8(value interface{}) int8 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < math.MinInt8 {
			return math.MinInt8
		}
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case int8:
		return x
	case int16:
		if x < math.MinInt8 {
			return math.MinInt8
		}
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case int32:
		if x < math.MinInt8 {
			return math.MinInt8
		}
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case int64:
		if x < math.MinInt8 {
			return math.MinInt8
		}
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case uint:
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case uint8:
		return int8(x)
	case uint16:
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case uint32:
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case uint64:
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case float32:
		if x < math.MinInt8 {
			return math.MinInt8
		}
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case float64:
		if x < math.MinInt8 {
			return math.MinInt8
		}
		if x > math.MaxInt8 {
			return math.MaxInt8
		}
		return int8(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 64)
		if err == nil {
			return Int8(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 64)
			if err == nil {
				return Int8(floatResult)
			}
		}
	case []byte:
		return Int8(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Int8(v.Float())
			} else if unsigned {
				return Int8(v.Uint())
			} else {
				return Int8(v.Int())
			}
		}
	}
	return 0
}

// Int16 将值转换成int16
func Int16(value interface{}) int16 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < math.MinInt16 {
			return math.MinInt16
		}
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case int8:
		return int16(x)
	case int16:
		return x
	case int32:
		if x < math.MinInt16 {
			return math.MinInt16
		}
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case int64:
		if x < math.MinInt16 {
			return math.MinInt16
		}
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case uint:
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case uint8:
		return int16(x)
	case uint16:
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case uint32:
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case uint64:
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case float32:
		if x < math.MinInt16 {
			return math.MinInt16
		}
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case float64:
		if x < math.MinInt16 {
			return math.MinInt16
		}
		if x > math.MaxInt16 {
			return math.MaxInt16
		}
		return int16(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 64)
		if err == nil {
			return Int16(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 64)
			if err == nil {
				return Int16(floatResult)
			}
		}
	case []byte:
		return Int16(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Int16(v.Float())
			} else if unsigned {
				return Int16(v.Uint())
			} else {
				return Int16(v.Int())
			}
		}
	}
	return 0
}

// Int32 将值转换成int32
func Int32(value interface{}) int32 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < math.MinInt32 {
			return math.MinInt32
		}
		if x > math.MaxInt32 {
			return math.MaxInt32
		}
		return int32(x)
	case int8:
		return int32(x)
	case int16:
		return int32(x)
	case int32:
		return x
	case int64:
		if x < math.MinInt32 {
			return math.MinInt32
		}
		if x > math.MaxInt32 {
			return math.MaxInt32
		}
		return int32(x)
	case uint:
		return int32(x)
	case uint8:
		return int32(x)
	case uint16:
		return int32(x)
	case uint32:
		if x > math.MaxInt32 {
			return math.MaxInt32
		}
		return int32(x)
	case uint64:
		if x > math.MaxInt32 {
			return math.MaxInt32
		}
		return int32(x)
	case float32:
		if x < math.MinInt32 {
			return math.MinInt32
		}
		if x > math.MaxInt32 {
			return math.MaxInt32
		}
		return int32(x)
	case float64:
		if x < math.MinInt32 {
			return math.MinInt32
		}
		if x > math.MaxInt32 {
			return math.MaxInt32
		}
		return int32(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 32)
		if err == nil {
			return Int32(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 32)
			if err == nil {
				return Int32(floatResult)
			}
		}
	case []byte:
		return Int32(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Int32(v.Float())
			} else if unsigned {
				return Int32(v.Uint())
			} else {
				return Int32(v.Int())
			}
		}
	}
	return 0
}

// Int64 将值转换成int64
func Int64(value interface{}) int64 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	case uint:
		return int64(x)
	case uint8:
		return int64(x)
	case uint16:
		return int64(x)
	case uint32:
		return int64(x)
	case uint64:
		if x > math.MaxInt64 {
			return math.MaxInt64
		}
		return int64(x)
	case float32:
		if x > math.MaxInt64 {
			return math.MaxInt64
		}
		return int64(x)
	case float64:
		if x < math.MinInt64 {
			return math.MinInt64
		}
		if x > math.MaxInt64 {
			return math.MaxInt64
		}
		return int64(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 64)
		if err == nil {
			return result
		} else {
			floatResult, err := strconv.ParseFloat(x, 64)
			if err == nil {
				return Int64(floatResult)
			}
		}
	case []byte:
		return Int64(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Int64(v.Float())
			} else if unsigned {
				return Int64(v.Uint())
			} else {
				return v.Int()
			}
		}
	}
	return 0
}

// Uint 将值转换成uint
func Uint(value interface{}) uint {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < 0 {
			return 0
		}
		return uint(x)
	case int8:
		if x < 0 {
			return 0
		}
		return uint(x)
	case int16:
		if x < 0 {
			return 0
		}
		return uint(x)
	case int32:
		if x < 0 {
			return 0
		}
		return uint(x)
	case int64:
		if x < 0 {
			return 0
		}
		return uint(x)
	case uint:
		return x
	case uint8:
		return uint(x)
	case uint16:
		return uint(x)
	case uint32:
		return uint(x)
	case uint64:
		return uint(x)
	case float32:
		if x < 0 {
			return 0
		}
		return uint(x)
	case float64:
		if x < 0 {
			return 0
		}
		return uint(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 32)
		if err == nil {
			return Uint(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 32)
			if err == nil {
				return Uint(floatResult)
			}
		}
	case []byte:
		return Uint(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Uint(v.Float())
			} else if unsigned {
				return Uint(v.Uint())
			} else {
				return Uint(v.Int())
			}
		}
	}
	return 0
}

// Uint8 将值转换成uint8
func Uint8(value interface{}) uint8 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case int8:
		if x < 0 {
			return 0
		}
		return uint8(x)
	case int16:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case int32:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case int64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case uint:
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case uint8:
		return x
	case uint16:
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case uint32:
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case uint64:
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case float32:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case float64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint8 {
			return math.MaxUint8
		}
		return uint8(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 32)
		if err == nil {
			return Uint8(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 32)
			if err == nil {
				return Uint8(floatResult)
			}
		}
	case []byte:
		return Uint8(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Uint8(v.Float())
			} else if unsigned {
				return Uint8(v.Uint())
			} else {
				return Uint8(v.Int())
			}
		}
	}
	return 0
}

// Uint16 将值转换成uint16
func Uint16(value interface{}) uint16 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case int8:
		if x < 0 {
			return 0
		}
		return uint16(x)
	case int16:
		if x < 0 {
			return 0
		}
		return uint16(x)
	case int32:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case int64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case uint:
		return uint16(x)
	case uint8:
		return uint16(x)
	case uint16:
		return x
	case uint32:
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case uint64:
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case float32:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case float64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint16 {
			return math.MaxUint16
		}
		return uint16(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 32)
		if err == nil {
			return Uint16(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 32)
			if err == nil {
				return Uint16(floatResult)
			}
		}
	case []byte:
		return Uint16(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Uint16(v.Float())
			} else if unsigned {
				return Uint16(v.Uint())
			} else {
				return Uint16(v.Int())
			}
		}
	}
	return 0
}

// Uint32 将值转换成uint32
func Uint32(value interface{}) uint32 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < 0 {
			return 0
		}
		return uint32(x)
	case int8:
		if x < 0 {
			return 0
		}
		return uint32(x)
	case int16:
		if x < 0 {
			return 0
		}
		return uint32(x)
	case int32:
		if x < 0 {
			return 0
		}
		return uint32(x)
	case int64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint32 {
			return math.MaxUint32
		}
		return uint32(x)
	case uint:
		return uint32(x)
	case uint8:
		return uint32(x)
	case uint16:
		return uint32(x)
	case uint32:
		return x
	case uint64:
		if x > math.MaxUint32 {
			return math.MaxUint32
		}
		return uint32(x)
	case float32:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint32 {
			return math.MaxUint32
		}
		return uint32(x)
	case float64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint32 {
			return math.MaxUint32
		}
		return uint32(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 32)
		if err == nil {
			return Uint32(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 32)
			if err == nil {
				return Uint32(floatResult)
			}
		}
	case []byte:
		return Uint32(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Uint32(v.Float())
			} else if unsigned {
				return Uint32(v.Uint())
			} else {
				return Uint32(v.Int())
			}
		}
	}
	return 0
}

// Uint64 将值转换成uint64
func Uint64(value interface{}) uint64 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int8:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int16:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int32:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int64:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case uint:
		return uint64(x)
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	case float32:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case float64:
		if x < 0 {
			return 0
		}
		if x > math.MaxUint64 {
			return math.MaxUint64
		}
		return uint64(x)
	case string:
		var result, err = strconv.ParseInt(x, 10, 64)
		if err == nil {
			return Uint64(result)
		} else {
			floatResult, err := strconv.ParseFloat(x, 64)
			if err == nil {
				return Uint64(floatResult)
			}
		}
	case []byte:
		return Uint64(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Uint64(v.Float())
			} else if unsigned {
				return v.Uint()
			} else {
				return Uint64(v.Int())
			}
		}
	}
	return 0
}

// Float64 将值转换成float64
func Float64(value interface{}) float64 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		return float64(x)
	case int8:
		return float64(x)
	case int16:
		return float64(x)
	case int32:
		return float64(x)
	case int64:
		return float64(x)
	case uint:
		return float64(x)
	case uint8:
		return float64(x)
	case uint16:
		return float64(x)
	case uint32:
		return float64(x)
	case uint64:
		return float64(x)
	case float32:
		return float64(x)
	case float64:
		return x
	case string:
		floatResult, err := strconv.ParseFloat(x, 64)
		if err == nil {
			return floatResult
		}
	case []byte:
		return Float64(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return v.Float()
			} else if unsigned {
				return Float64(v.Uint())
			} else {
				return Float64(v.Int())
			}
		}
	}
	return 0
}

// Float32 将值转换成float32
func Float32(value interface{}) float32 {
	if value == nil {
		return 0
	}

	switch x := value.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		return float32(x)
	case int8:
		return float32(x)
	case int16:
		return float32(x)
	case int32:
		return float32(x)
	case int64:
		if float64(x) > math.MaxFloat32 {
			return math.MaxFloat32
		}
		return float32(x)
	case uint:
		return float32(x)
	case uint8:
		return float32(x)
	case uint16:
		return float32(x)
	case uint32:
		return float32(x)
	case uint64:
		return float32(x)
	case float32:
		return x
	case float64:
		if x > math.MaxFloat32 {
			return math.MaxFloat32
		}
		return float32(x)
	case string:
		floatResult, err := strconv.ParseFloat(x, 32)
		if err == nil {
			return Float32(floatResult)
		}
	case []byte:
		return Float32(string(x))
	default:
		var v = reflect.ValueOf(value)
		unsigned, isFloat, ok := IsNumberKind(v.Kind())
		if ok {
			if isFloat {
				return Float32(v.Float())
			} else if unsigned {
				return Float32(v.Uint())
			} else {
				return Float32(v.Int())
			}
		}
	}
	return 0
}

// Bool 将值转换成bool类型
func Bool(value interface{}) bool {
	if value == nil {
		return false
	}

	var kind = reflect.TypeOf(value).Kind()
	switch kind {
	case reflect.Bool:
		return value.(bool)
	}
	return Int64(value) > 0
}

// String 将值转换成字符串
func String(value interface{}) string {
	if value == nil {
		return ""
	}
	switch x := value.(type) {
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case int8, int16, int32, uint, uint8, uint16, uint32:
		return fmt.Sprintf("%d", value)
	case []byte:
		return string(value.([]byte))
	case string:
		return x
	}
	return fmt.Sprintf("%#v", value)
}

// Compare 比较两个值大小
func Compare(value1 interface{}, value2 interface{}) bool {
	if value1 == nil {
		return false
	}

	switch value1 := value1.(type) {
	case bool:
		return Int(value1) > Int(value2)
	case int:
		return Int(value1) > Int(value2)
	case int8:
		return Int8(value1) > Int8(value2)
	case int16:
		return Int16(value1) > Int16(value2)
	case int32:
		return Int32(value1) > Int32(value2)
	case int64:
		return Int64(value1) > Int64(value2)
	case uint:
		return Uint(value1) > Uint(value2)
	case uint8:
		return Uint8(value1) > Uint8(value2)
	case uint16:
		return Uint16(value1) > Uint16(value2)
	case uint32:
		return Uint32(value1) > Uint32(value2)
	case uint64:
		return Uint64(value1) > Uint64(value2)
	case float32:
		return Float32(value1) > Float32(value2)
	case float64:
		return Float64(value1) > Float64(value2)
	case string:
		return String(value1) > String(value2)
	}
	return String(value1) > String(value2)
}

// IsNumber 判断是否为数字
func IsNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	}
	return false
}

// IsNumberKind 判断是否为数字类型
func IsNumberKind(kind reflect.Kind) (unsigned bool, isFloat bool, ok bool) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		unsigned = false
		ok = true
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		unsigned = true
		ok = true
		return
	case reflect.Float32, reflect.Float64:
		isFloat = true
		ok = true
		return
	}
	return
}

// IsInteger 判断是否为整形数字
func IsInteger(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

// IsFloat 判断是否为浮点型数字
func IsFloat(value interface{}) bool {
	switch value.(type) {
	case float32, float64:
		return true
	}
	return false
}

// IsSlice 判断是否为Slice
func IsSlice(value interface{}) bool {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

// IsMap 判断是否为Map
func IsMap(value interface{}) bool {
	if value == nil {
		return false
	}
	return reflect.TypeOf(value).Kind() == reflect.Map
}

// IsNil 判断是否为nil
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}

	return reflect.ValueOf(value).IsNil()
}

// Slice 转换Slice类型
func Slice(fromSlice interface{}, toSliceType reflect.Type) (interface{}, error) {
	if fromSlice == nil {
		return nil, errors.New("'fromSlice' should not be nil")
	}

	fromValue := reflect.ValueOf(fromSlice)
	if fromValue.Kind() != reflect.Slice {
		return nil, errors.New("'fromSlice' should be slice")
	}

	if toSliceType.Kind() != reflect.Slice {
		return nil, errors.New("'toSliceType' should be slice")
	}

	v := reflect.Indirect(reflect.New(toSliceType))
	count := fromValue.Len()
	toElemKind := toSliceType.Elem().Kind()
	for i := 0; i < count; i++ {
		elem := fromValue.Index(i)
		elemVar := elem.Interface()
		switch toElemKind {
		case reflect.Int:
			v = reflect.Append(v, reflect.ValueOf(Int(elemVar)))
		case reflect.Int8:
			v = reflect.Append(v, reflect.ValueOf(Int8(elemVar)))
		case reflect.Int16:
			v = reflect.Append(v, reflect.ValueOf(Int16(elemVar)))
		case reflect.Int32:
			v = reflect.Append(v, reflect.ValueOf(Int32(elemVar)))
		case reflect.Int64:
			v = reflect.Append(v, reflect.ValueOf(Int64(elemVar)))
		case reflect.Uint:
			v = reflect.Append(v, reflect.ValueOf(Uint(elemVar)))
		case reflect.Uint8:
			v = reflect.Append(v, reflect.ValueOf(Uint8(elemVar)))
		case reflect.Uint16:
			v = reflect.Append(v, reflect.ValueOf(Uint16(elemVar)))
		case reflect.Uint32:
			v = reflect.Append(v, reflect.ValueOf(Uint32(elemVar)))
		case reflect.Uint64:
			v = reflect.Append(v, reflect.ValueOf(Uint64(elemVar)))
		case reflect.Bool:
			v = reflect.Append(v, reflect.ValueOf(Bool(elemVar)))
		case reflect.Float32:
			v = reflect.Append(v, reflect.ValueOf(Float32(elemVar)))
		case reflect.Float64:
			v = reflect.Append(v, reflect.ValueOf(Float64(elemVar)))
		case reflect.String:
			v = reflect.Append(v, reflect.ValueOf(String(elemVar)))
		}
	}
	return v.Interface(), nil
}
