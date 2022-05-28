package utils

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// RandDecimalByWeight 根据小数点位数权重随机
func RandDecimalByWeight(num float64, w string) (float64, error) {
	weights := FormatWeight(w)

	weight, err := weights.Rand()
	if err != nil {
		return 0, err
	}

	decimal, err := strconv.Atoi(weight)
	if err != nil {
		return 0, err
	}

	return Round(num, decimal), nil
}

// RandFloat64 随机数
func RandFloat64(from float64, to float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(to-from) + from
}

// RandSleepTime 随机暂停时间
func RandSleepTime(milli int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(milli)
}

// RandDecimal 随机数
func RandDecimal(from decimal.Decimal, to decimal.Decimal) decimal.Decimal {
	rand.Seed(time.Now().UnixNano())
	return decimal.NewFromFloat(rand.Float64()).Mul(to.Sub(from)).Add(from)
}
