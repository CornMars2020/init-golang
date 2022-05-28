package utils

import (
	"github.com/shopspring/decimal"
)

// BenchmarkPriceItem 对标价格数据, [5]decimal.Decimal{askPrice, bidPrice, midPrice, wmidPrice, weight}
type BenchmarkPriceItem = [5]decimal.Decimal

// BenchmarkPrices 对标价格列表
type BenchmarkPrices []BenchmarkPriceItem

func (bp BenchmarkPrices) Len() int                     { return len(bp) }
func (bp BenchmarkPrices) Less(i, j int) bool           { return bp[i][2].LessThan(bp[j][2]) }
func (bp BenchmarkPrices) Swap(i, j int)                { bp[i], bp[j] = bp[j], bp[i] }
func (bp BenchmarkPrices) MidIndex() int                { return bp.Len() / 2 }
func (bp BenchmarkPrices) MidValue() BenchmarkPriceItem { return bp[bp.MidIndex()] }

// GetAvgPrice 获取列表价格算数平均价
func GetAvgPrice(prices BenchmarkPrices) (ask1 decimal.Decimal, bid1 decimal.Decimal) {
	var weight decimal.Decimal
	var task, tbid decimal.Decimal
	for _, price := range prices {
		if price[0].LessThanOrEqual(decimal.Zero) || price[1].LessThanOrEqual(decimal.Zero) {
			continue
		}
		weight = weight.Add(decimal.NewFromInt(1))
		task = task.Add(price[0])
		tbid = tbid.Add(price[1])
	}

	if weight.LessThanOrEqual(decimal.Zero) {
		return
	}

	ask1 = task.Div(weight)
	bid1 = tbid.Div(weight)
	return
}

// GetWAvgPrice 获取列表价格加权平均价
func GetWAvgPrice(prices BenchmarkPrices) (ask1 decimal.Decimal, bid1 decimal.Decimal) {
	var weight decimal.Decimal
	var task, tbid decimal.Decimal
	for _, price := range prices {
		if price[0].LessThanOrEqual(decimal.Zero) || price[1].LessThanOrEqual(decimal.Zero) {
			continue
		}
		weight = weight.Add(price[4])
		task = price[0].Mul(price[4]).Add(task)
		tbid = price[1].Mul(price[4]).Add(tbid)
	}

	if weight.LessThanOrEqual(decimal.Zero) {
		return
	}

	ask1 = task.Div(weight)
	bid1 = tbid.Div(weight)
	return
}

// GetWorstPrice 获取列表价格最差价
func GetWorstPrice(prices BenchmarkPrices) (ask1 decimal.Decimal, bid1 decimal.Decimal) {
	for _, price := range prices {
		if price[0].LessThanOrEqual(decimal.Zero) || price[1].LessThanOrEqual(decimal.Zero) {
			continue
		}

		if ask1.IsZero() || ask1.LessThan(price[0]) {
			ask1 = price[0]
		}

		if bid1.IsZero() || bid1.LessThan(price[1]) {
			bid1 = price[1]
		}
	}
	return
}

// GetBestPrice 获取列表价格最优价
func GetBestPrice(prices BenchmarkPrices) (ask1 decimal.Decimal, bid1 decimal.Decimal) {
	for _, price := range prices {
		if price[0].LessThanOrEqual(decimal.Zero) || price[1].LessThanOrEqual(decimal.Zero) {
			continue
		}

		if ask1.IsZero() || ask1.GreaterThan(price[0]) {
			ask1 = price[0]
		}

		if bid1.IsZero() || bid1.LessThan(price[1]) {
			bid1 = price[1]
		}
	}
	return
}

// GetMidPrice 基于中间价向外扩展铺单
func GetMidPrice(prices BenchmarkPrices) (ask1 decimal.Decimal, bid1 decimal.Decimal) {
	var weight decimal.Decimal
	var task, tbid decimal.Decimal
	for _, price := range prices {
		if price[2].LessThanOrEqual(decimal.Zero) {
			continue
		}
		weight = weight.Add(decimal.NewFromInt(1))
		task = task.Add(price[2])
		tbid = tbid.Add(price[2])
	}

	if weight.LessThanOrEqual(decimal.Zero) {
		return
	}

	ask1 = task.Div(weight)
	bid1 = tbid.Div(weight)
	return
}

// GetWMidPrice 基于加权中间价向外扩展铺单
func GetWMidPrice(prices BenchmarkPrices) (ask1 decimal.Decimal, bid1 decimal.Decimal) {
	var weight decimal.Decimal
	var task, tbid decimal.Decimal
	for _, price := range prices {
		if price[2].LessThanOrEqual(decimal.Zero) {
			continue
		}
		weight = weight.Add(decimal.NewFromInt(1))
		task = task.Add(price[3])
		tbid = tbid.Add(price[3])
	}

	if weight.LessThanOrEqual(decimal.Zero) {
		return
	}

	ask1 = task.Div(weight)
	bid1 = tbid.Div(weight)
	return
}

// GetMarkPrice 获取盘口中间价
func GetMarkPrice(ask decimal.Decimal, bid decimal.Decimal) decimal.Decimal {
	mid := ask.Add(bid).Div(decimal.NewFromInt(2))
	return mid
}

// GetAskBidSpread 获取盘口差
func GetAskBidSpread(ask decimal.Decimal, bid decimal.Decimal) decimal.Decimal {
	invalidValue := decimal.NewFromInt(-1)
	if bid.IsZero() || ask.LessThan(bid) || ask.IsNegative() || bid.IsNegative() {
		return invalidValue
	}

	spread := ask.Sub(bid)
	return spread
}

// GetAskBidPercent 获取盘口差比例
func GetAskBidPercent(ask decimal.Decimal, bid decimal.Decimal) decimal.Decimal {
	invalidValue := decimal.NewFromInt(-1)
	if bid.IsZero() || ask.LessThan(bid) || ask.IsNegative() || bid.IsNegative() {
		return invalidValue
	}

	spread := ask.Sub(bid)
	percent := spread.Div(bid)
	return percent
}
