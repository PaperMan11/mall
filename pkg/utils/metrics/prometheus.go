package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
)

const (
	namespace = "ziqing"
	subsystem = "mall"
)

// metricsRequestsTotal metrics for request total 计数器（Counter）
var metricsRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_total",
		Help:      "request(ms) total",
	},
	[]string{"method", "path"},
)

// metricsRequestsCost metrics for requests cost 累积直方图（Histogram）
var metricsRequestsCost = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_cost",
		Help:      "request(ms) cost milliseconds",
	},
	[]string{"method", "path", "success", "http_code", "cost_milliseconds"},
)

func init() {
	prometheus.MustRegister(metricsRequestsTotal, metricsRequestsCost)
}

func RecordMetrics(method, path string, success bool, httpCode int, costSeconds float64) {
	metricsRequestsTotal.With(prometheus.Labels{
		"method": method,
		"path":   path,
	}).Inc()

	metricsRequestsCost.With(prometheus.Labels{
		"method":            method,
		"path":              path,
		"success":           cast.ToString(success),
		"http_code":         cast.ToString(httpCode),
		"cost_milliseconds": cast.ToString(costSeconds * 1000),
	}).Observe(costSeconds)
}
