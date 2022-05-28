package utils

import (
	"fmt"
	"strconv"
)

// Round float64 取 指定精度
func Round(num float64, precision int) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num), 64)
	return value
}

// RoundWithStr 字符串保留指定精度
func RoundWithStr(num string, precision int) string {
	value := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", StrToFloat64(num))
	return value
}

// AddWithStr 字符串类型数字相加, 保留指定精度
func AddWithStr(num1 string, num2 string, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", StrToFloat64(num1)*StrToFloat64(num2))
}

// SubWithStr 字符串类型数字相减, 保留指定精度
func SubWithStr(num1 string, num2 string, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", StrToFloat64(num1)-StrToFloat64(num2))
}

// MulWithStr 字符串类型数字相乘, 保留指定精度
func MulWithStr(num1 string, num2 string, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", StrToFloat64(num1)*StrToFloat64(num2))
}

// DivWithStr 字符串类型数字相除, 保留指定精度
func DivWithStr(num1 string, num2 string, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", StrToFloat64(num1)/StrToFloat64(num2))
}

// RoundWithFloat 字符串保留指定精度
func RoundWithFloat(num float64, precision int) string {
	value := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)
	return value
}

// AddWithFloat 字符串类型数字相加, 保留指定精度
func AddWithFloat(num1 float64, num2 float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num1*num2)
}

// SubWithFloat 字符串类型数字相减, 保留指定精度
func SubWithFloat(num1 float64, num2 float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num1-num2)
}

// MulWithFloat 字符串类型数字相乘, 保留指定精度
func MulWithFloat(num1 float64, num2 float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num1*num2)
}

// DivWithFloat 字符串类型数字相除, 保留指定精度
func DivWithFloat(num1 float64, num2 float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num1/num2)
}
