package ycmonitoringgo

import (
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type Rate struct {
	name   string
	labels []string

	metrics map[string]*rateMetric
	mu      sync.RWMutex
}

type rateMetric struct {
	Value       atomic.Float64
	LabelValues []string
}

func NewRate(name string, registry *Registry, labels ...string) *Rate {
	r := &Rate{
		name:   name,
		labels: labels,

		metrics: make(map[string]*rateMetric),
	}

	registry.Add(r)
	return r
}

func (s *Rate) Inc(values ...string) {
	s.Add(1, values...)
}

func (s *Rate) Add(delta float64, values ...string) {
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
		metric = &rateMetric{
			LabelValues: values,
		}
		s.metrics[tagKey] = metric
	}

	metric.Value.Add(delta)
	s.mu.Unlock()
}

func (s *Rate) Reset(values ...string) {
	tagKey := strings.Join(values, ",")

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.metrics, tagKey)
}

func (s *Rate) Name() string {
	return s.name
}

func (s *Rate) GetMetrics() []metric {
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
			Type:   TYPE_RATE,
			Value:  m.Value.Load(),
		})
	}

	return result
}
