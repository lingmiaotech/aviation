package prom

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	"github.com/lingmiaotech/tonic/configs"
	"github.com/lingmiaotech/tonic/logging"
	"github.com/lingmiaotech/tonic/prom/pool"
)

type InstanceClass struct {
	AppName     string
	MetricsPool *pool.Pool
}

type Timer struct {
	Start time.Time
}

var Instance InstanceClass

func Increment(bucket string) {
	name, ok := getBucket(bucket)
	if !ok {
		logging.GetDefaultLogger().Infof("[PROM] key=%s count=1", bucket)
		return
	}
	Instance.MetricsPool.GetCounter(name).Inc()
}

// Timing takes bucket name and delta in milliseconds
func Timing(bucket string, delta int) {
	name, ok := getBucket(bucket)
	if !ok {
		logging.GetDefaultLogger().Infof("[PROM] key=%s time_delta=%d(ms)", bucket, delta)
		return
	}
	Instance.MetricsPool.GetHistogram(name).Observe(float64(delta) / 10e5)
}

// Count increments bucket name by n
func Count(bucket string, n int) {
	name, ok := getBucket(bucket)
	if !ok {
		logging.GetDefaultLogger().Infof("[PROM] key=%v count=%d", bucket, n)
		return
	}
	Instance.MetricsPool.GetCounter(name).Add(float64(n))
}

// Gauge set bucket name to n
func Gauge(bucket string, n int) {
	name, ok := getBucket(bucket)
	if !ok {
		logging.GetDefaultLogger().Infof("[PROM] key=%v gauge=%d", bucket, n)
		return
	}
	Instance.MetricsPool.GetGauge(name).Set(float64(n))
}

func isLegit(bucket string) bool {
	return model.MetricNameRE.MatchString(bucket)
}

func getBucket(bucket string) (string, bool) {
	name := strings.ReplaceAll(fmt.Sprintf("%s.%s", Instance.AppName, bucket), ".", ":")
	return name, isLegit(name)
}

// NewTimer returns a Timer
func NewTimer() Timer {
	return Timer{Start: time.Now()}
}

// Send sends the time elapsed since the creation of the Timing.
func (t Timer) Send(bucket string) {
	Timing(bucket, int(t.Duration().Nanoseconds()/1000))
}

// Duration returns the time elapsed since the creation of the Timing.
func (t Timer) Duration() time.Duration {
	return time.Since(t.Start)
}

func NewCustomTimer(t time.Time) Timer {
	return Timer{Start: t}
}

func InitProm() error {
	Instance.AppName = configs.GetString("app_name")
	Instance.MetricsPool = &pool.Pool{
		Counters:   make(map[string]prometheus.Counter),
		Gauges:     make(map[string]prometheus.Gauge),
		Histograms: make(map[string]prometheus.Histogram),
		Summaries:  make(map[string]prometheus.Summary),
	}
	return nil
}
