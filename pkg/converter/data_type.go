package converter

import (
	"strconv"
	"time"
)

// StringToFloat64 returns the float64 value from the input string
func StringToFloat64(input string) float64 {
	result, _ := strconv.ParseFloat(input, 64)
	return result
}

// StringToFloat32 returns the float32 value from the input string
func StringToFloat32(input string) float32 {
	result, _ := strconv.ParseFloat(input, 32)
	return float32(result)
}

// StringToInt returns the int value from the input string
func StringToInt(input string) int {
	result, _ := strconv.Atoi(input)
	return result
}

// StringToInt8 returns the int8 value from the input string
func StringToInt8(input string) int8 {
	result, _ := strconv.ParseInt(input, 10, 8)
	return int8(result)
}

// StringToInt16 returns the int16 value from the input string
func StringToInt16(input string) int16 {
	result, _ := strconv.ParseInt(input, 10, 16)
	return int16(result)
}

// StringToInt32 returns the int32 value from the input string
func StringToInt32(input string) int32 {
	result, _ := strconv.ParseInt(input, 10, 32)
	return int32(result)
}

// StringToInt64 returns the int64 value from the input string
func StringToInt64(input string) int64 {
	result, _ := strconv.ParseInt(input, 10, 64)
	return result
}

// StringToUint returns the uint value from the input string
func StringToUint(input string) uint {
	result, _ := strconv.ParseUint(input, 10, 0)
	return uint(result)
}

// StringToUint8 returns the uint8 value from the input string
func StringToUint8(input string) uint8 {
	result, _ := strconv.ParseUint(input, 10, 8)
	return uint8(result)
}

// StringToUint16 returns the uint16 value from the input string
func StringToUint16(input string) uint16 {
	result, _ := strconv.ParseUint(input, 10, 16)
	return uint16(result)
}

// StringToUint32 returns the uint32 value from the input string
func StringToUint32(input string) uint32 {
	result, _ := strconv.ParseUint(input, 10, 32)
	return uint32(result)
}

// StringToUint64 returns the uint64 value from the input string
func StringToUint64(input string) uint64 {
	result, _ := strconv.ParseUint(input, 10, 64)
	return result
}

// StringToBool returns the bool value from the input string
func StringToBool(input string) bool {
	result, _ := strconv.ParseBool(input)
	return result
}

// StringToTimeValue convert string to time based on the provided format. zero time return if invalid
func StringToTimeValue(val, format string) time.Time {
	timeVal, _ := time.Parse(format, val)

	return timeVal
}
