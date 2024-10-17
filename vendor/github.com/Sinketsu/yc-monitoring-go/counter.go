package ycmonitoringgo

import (
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type Counter struct {
	name   string
	labels []string

	metrics map[string]*counterMetric
	mu      sync.RWMutex
}

type counterMetric struct {
	Value       atomic.Int64
	LabelValues []string
}

func NewCounter(name string, registry *Registry, labels ...string) *Counter {
	c := &Counter{
		name:   name,
		labels: labels,

		metrics: make(map[string]*counterMetric),
	}

	registry.Add(c)
	return c
}

func (s *Counter) Inc(values ...string) {
	s.Add(1, values...)
}

func (s *Counter) Add(delta int64, values ...string) {
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
		metric = &counterMetric{
			LabelValues: values,
		}
		s.metrics[tagKey] = metric
	}

	metric.Value.Add(delta)
	s.mu.Unlock()
}

func (s *Counter) Reset(values ...string) {
	tagKey := strings.Join(values, ",")

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.metrics, tagKey)
}

func (s *Counter) Name() string {
	return s.name
}

func (s *Counter) GetMetrics() []metric {
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
			Type:   TYPE_COUNTER,
			Value:  m.Value.Load(),
		})
	}

	return result
}
