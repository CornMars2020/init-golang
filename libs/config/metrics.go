package config

import (
	"fmt"
	"strings"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"
)

// Metrics
type Metrics struct {
	APICounter   *kitprometheus.Counter // API 计数器
	APISummary   *kitprometheus.Summary // API 延时统计
	SvcGauge     *kitprometheus.Gauge   // 服务状态
	ErrorDetails bool
}

func NewMetrics(ns string, sys string, details bool) *Metrics {
	fieldKeys := []string{"svc", "api", "error"}
	apiCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: ns,
		Subsystem: sys,
		Name:      "api_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	apiSummary := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: ns,
		Subsystem: sys,
		Name:      "api_latency_ms",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	fieldKeys = []string{"svc", "name", "type"}
	svcGauge := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: ns,
		Subsystem: sys,
		Name:      "service_status",
		Help:      "Service status in different view point.",
	}, fieldKeys)

	metrics := &Metrics{
		APICounter:   apiCount,
		APISummary:   apiSummary,
		SvcGauge:     svcGauge,
		ErrorDetails: details,
	}

	return metrics
}

func getLastError(err error) string {
	msg := err.Error()

	msgs := strings.Split(msg, " ")
	length := len(msgs)

	lst := msgs[length-1]

	return lst
}

func (metrics *Metrics) AddApiCounter(svc string, api string, err error) {
	var lvs []string
	if metrics.ErrorDetails && err != nil {
		lvs = []string{"svc", svc, "api", api, "error", getLastError(err)}
	} else {
		lvs = []string{"svc", svc, "api", api, "error", fmt.Sprint(err != nil)}
	}
	metrics.APICounter.With(lvs...).Add(1)
}

func (metrics *Metrics) SetApiSummary(svc string, api string, err error, begin time.Time) {
	var lvs []string
	if metrics.ErrorDetails && err != nil {
		lvs = []string{"svc", svc, "api", api, "error", getLastError(err)}
	} else {
		lvs = []string{"svc", svc, "api", api, "error", fmt.Sprint(err != nil)}
	}
	metrics.APISummary.With(lvs...).Observe(time.Since(begin).Seconds())
}

func (metrics *Metrics) SetSvcValue(svc string, name string, tp string, value float64) {
	lvs := []string{"svc", svc, "name", name, "type", tp}
	metrics.SvcGauge.With(lvs...).Set(value)
}

func (metrics *Metrics) SetSvcInt16(svc string, name string, tp string, value int16) {
	v := float64(value)
	metrics.SetSvcValue(svc, name, tp, v)
}
func (metrics *Metrics) SetSvcInt64(svc string, name string, tp string, value int64) {
	v := float64(value)
	metrics.SetSvcValue(svc, name, tp, v)
}
func (metrics *Metrics) SetSvcFloat64(svc string, name string, tp string, value float64) {
	metrics.SetSvcValue(svc, name, tp, value)
}
func (metrics *Metrics) SetSvcDecimal(svc string, name string, tp string, value decimal.Decimal) {
	v, _ := value.Float64()
	metrics.SetSvcValue(svc, name, tp, v)
}
