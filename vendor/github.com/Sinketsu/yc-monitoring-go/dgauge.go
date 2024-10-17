package ycmonitoringgo

import (
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type DGauge struct {
	name   string
	labels []string

	metrics map[string]*dgaugeMetric
	mu      sync.RWMutex
}

type dgaugeMetric struct {
	Value       atomic.Float64
	LabelValues []string
}

func NewDGauge(name string, registry *Registry, labels ...string) *DGauge {
	dg := &DGauge{
		name:   name,
		labels: labels,

		metrics: make(map[string]*dgaugeMetric),
	}

	registry.Add(dg)
	return dg
}

func (s *DGauge) Set(value float64, values ...string) {
	if len(values) != len(s.labels) {
		return
	}
	tagKey := strings.Join(values, ",")

	s.mu.RLock()
	metric, ok := s.metrics[tagKey]
	s.mu.RUnlock()

	if ok {
		metric.Value.Store(value)
		return
	}

	s.mu.Lock()
	metric, ok = s.metrics[tagKey]
	if !ok {
		metric = &dgaugeMetric{
			LabelValues: values,
		}
		s.metrics[tagKey] = metric
	}

	metric.Value.Store(value)
	s.mu.Unlock()
}

func (s *DGauge) Add(delta float64, values ...string) {
	if len(values) != len(s.labels) {
		return
	}
	tagKey := strings.Join(values, ",")

	s.mu.RLock()
	metric, ok := s.metrics[tagKey]
	s.mu.RUnlock()

	if ok {
		metric.Value.Add(delta)
		return
	}

	s.mu.Lock()
	metric, ok = s.metrics[tagKey]
	if !ok {
		metric = &dgaugeMetric{
			LabelValues: values,
		}
		s.metrics[tagKey] = metric
	}

	metric.Value.Add(delta)
	s.mu.Unlock()
}

func (s *DGauge) Reset(values ...string) {
	tagKey := strings.Join(values, ",")

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.metrics, tagKey)
}

func (s *DGauge) Name() string {
	return s.name
}

func (s *DGauge) GetMetrics() []metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]metric, 0, len(s.metrics))
	for _, m := range s.metrics {
		labels := make(map[string]string, len(s.labels))
		for i, name := range s.labels {
			labels[name] = m.LabelValues[i]
		}

		result = append(result, metric{
			Name:   s.name,
			Labels: labels,
			Type:   TYPE_DGAUGE,
			Value:  m.Value.Load(),
		})
	}

	return result
}
