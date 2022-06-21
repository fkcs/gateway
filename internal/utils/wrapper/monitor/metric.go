package monitor

import (
	"github.com/fkcs/gateway/internal/utils/wrapper/monitor/prometheus"
	"strconv"
	"time"
)

const serverNamespace = "http_server"

var (
	MetricServerReqDur = prometheus.NewHistogramVec(&prometheus.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(s).",
		Labels:    []string{"path"},
		Buckets:   []float64{0.01, 0.1, 0.5, 1, 2, 10, 30, 60},
	})

	MetricServerReqCodeTotal = prometheus.NewCounterVec(&prometheus.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "code"},
	})
)

func ObserveMetric(cost time.Duration, path string, code int) {
	MetricServerReqDur.Observe(cost.Seconds(), path)
	MetricServerReqCodeTotal.Inc(path, strconv.Itoa(code))
}
