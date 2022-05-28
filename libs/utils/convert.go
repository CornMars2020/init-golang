package utils

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// StrToFloat64 字符串转 float64
func StrToFloat64(num string) float64 {
	value, _ := strconv.ParseFloat(num, 64)
	return value
}

// Float64ToStr float64 转 字符串
func Float64ToStr(num float64, precision int) string {
	return strconv.FormatFloat(num, 'f', precision, 64)
}

// StrToInt64 字符串转 int64
func StrToInt64(num string) int64 { 
	value, _ := strconv.ParseInt(num, 10, 64)
	return value
}

// Int64ToStr int64 转字符串
func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

// StrToInt 字符串转 int
func StrToInt(num string, defaultValue int) int {
	if num == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(num)
	if err != nil {
		return defaultValue
	}
	return value
}

func DecimalToFixedDecimal(v decimal.Decimal, decimals int) decimal.Decimal {
	fixed := v.StringFixed(int32(decimals))
	value, _ := decimal.NewFromString(fixed)
	return value
}

func DecimalToFixed(v decimal.Decimal, decimals int) string {
	return v.StringFixed(int32(decimals))
}
