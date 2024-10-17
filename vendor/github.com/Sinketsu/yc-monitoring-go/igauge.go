package ycmonitoringgo

import (
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type IGauge struct {
	name   string
	labels []string

	metrics map[string]*igaugeMetric
	mu      sync.RWMutex
}

type igaugeMetric struct {
	Value       atomic.Int64
	LabelValues []string
}

func NewIGauge(name string, registry *Registry, labels ...string) *IGauge {
	ig := &IGauge{
		name:   name,
		labels: labels,

		metrics: make(map[string]*igaugeMetric),
	}

	registry.Add(ig)
	return ig
}

func (s *IGauge) Set(value int64, values ...string) {
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
		metric = &igaugeMetric{
			LabelValues: values,
		}
		s.metrics[tagKey] = metric
	}

	metric.Value.Store(value)
	s.mu.Unlock()
}

func (s *IGauge) Add(delta int64, values ...string) {
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
		metric = &igaugeMetric{
			LabelValues: values,
		}
		s.metrics[tagKey] = metric
	}

	metric.Value.Add(delta)
	s.mu.Unlock()
}

func (s *IGauge) Reset(values ...string) {
	tagKey := strings.Join(values, ",")

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.metrics, tagKey)
}

func (s *IGauge) Name() string {
	return s.name
}

func (s *IGauge) GetMetrics() []metric {
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
			Type:   TYPE_IGAUGE,
			Value:  m.Value.Load(),
		})
	}

	return result
}
