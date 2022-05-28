package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type weightItem struct {
	item   string
	weight int
}

// Weight 权重
type Weight struct {
	Total   int
	Weights []weightItem
}

// Rand 依据权重随机值
func (w *Weight) Rand() (string, error) {
	if len(w.Weights) <= 0 {
		return "", fmt.Errorf("invalid weights len")
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(w.Total)

	var total int = 0
	for _, kv := range w.Weights {
		total = total + kv.weight
		if id < total {
			return kv.item, nil
		}
	}

	return "", fmt.Errorf("invalid weights item")
}

// FormatWeight 字符串转换为权重值
func FormatWeight(weight string) *Weight {
	var items []string = strings.Split(weight, "|")
	var weights []weightItem
	var total int = 0
	for _, item := range items {
		var kv []string = strings.Split(item, ":")
		if len(kv) == 2 && kv[0] != "" && kv[1] != "" {

			v, err := strconv.Atoi(kv[1])
			if err != nil {
				continue
			}
			weights = append(weights, weightItem{item: kv[0], weight: v})
			total = total + v
		}
	}
	return &Weight{
		Total:   total,
		Weights: weights,
	}
}
