package pool

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Pool struct {
	CountersLock sync.RWMutex
	Counters     map[string]prometheus.Counter

	GaugesLock sync.RWMutex
	Gauges     map[string]prometheus.Gauge

	HistogramsLock sync.RWMutex
	Histograms     map[string]prometheus.Histogram

	SummariesLock sync.RWMutex
	Summaries     map[string]prometheus.Summary
}

func (p *Pool) GetCounter(name string) prometheus.Counter {
	p.CountersLock.RLock()
	v, ok := p.Counters[name]
	p.CountersLock.RUnlock()
	if ok {
		return v
	}
	return p.AddCounter(name)
}

func (p *Pool) GetGauge(name string) prometheus.Gauge {
	p.GaugesLock.RLock()
	v, ok := p.Gauges[name]
	p.GaugesLock.RUnlock()
	if ok {
		return v
	}
	return p.AddGauge(name)
}

func (p *Pool) GetHistogram(name string) prometheus.Histogram {
	p.HistogramsLock.RLock()
	v, ok := p.Histograms[name]
	p.HistogramsLock.RUnlock()
	if ok {
		return v
	}
	return p.AddHistogram(name)
}

func (p *Pool) GetSummary(name string) prometheus.Summary {
	p.SummariesLock.RLock()
	v, ok := p.Summaries[name]
	p.SummariesLock.RUnlock()
	if ok {
		return v
	}
	return p.AddSummary(name)
}

func (p *Pool) AddCounter(name string) prometheus.Counter {
	p.CountersLock.Lock()
	defer p.CountersLock.Unlock()
	opts := prometheus.CounterOpts{
		Name: name,
		Help: name + "help",
	}
	v := promauto.NewCounter(opts)
	p.Counters[name] = v
	return v
}

func (p *Pool) AddGauge(name string) prometheus.Gauge {
	p.GaugesLock.Lock()
	defer p.GaugesLock.Unlock()
	opts := prometheus.GaugeOpts{
		Name: name,
		Help: name + "help",
	}
	v := promauto.NewGauge(opts)
	p.Gauges[name] = v
	return v
}

func (p *Pool) AddSummary(name string) prometheus.Summary {
	p.SummariesLock.Lock()
	defer p.SummariesLock.Unlock()
	opts := prometheus.SummaryOpts{
		Name:       name,
		Help:       name + "help",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}
	v := promauto.NewSummary(opts)
	p.Summaries[name] = v
	return v
}

func (p *Pool) AddHistogram(name string) prometheus.Histogram {
	p.HistogramsLock.Lock()
	defer p.HistogramsLock.Unlock()
	opts := prometheus.HistogramOpts{
		Name:    name,
		Help:    name + "help",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.2, 0.3, 0.5, 0.8, 1.3, 2.1, 3.4, 5.5},
	}
	v := promauto.NewHistogram(opts)
	p.Histograms[name] = v
	return v
}
