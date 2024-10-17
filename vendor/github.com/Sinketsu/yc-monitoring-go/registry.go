package ycmonitoringgo

import (
	"fmt"
	"sync"
)

var (
	DefaultRegistry = NewRegistry()
)

type Registry struct {
	names   map[string]struct{}
	metrics []Metric
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		names: make(map[string]struct{}),
	}
}

func (r *Registry) Add(metric Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.names[metric.Name()]; ok {
		panic(fmt.Sprintf("multiple metric '%s' registration found", metric.Name()))
	}

	r.metrics = append(r.metrics, metric)
	r.names[metric.Name()] = struct{}{}
}

func (r *Registry) Range(f func(i int, m Metric)) {
	r.mu.RLock()

	for i, m := range r.metrics {
		f(i, m)
	}

	r.mu.RUnlock()
}
