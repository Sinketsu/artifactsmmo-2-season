package monitoring

import (
	"slices"
	"sync"
)

type metric struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
	Type   string            `json:"type"`
	Value  float64           `json:"value"`
}

type DGaugeMetric struct {
	Value       float64
	LabelValues []string
}

type DGauge struct {
	name   string
	labels []string

	metrics []DGaugeMetric
	mu      sync.RWMutex
}

func NewDGauge(name string, labels ...string) *DGauge {
	dg := &DGauge{
		name:   name,
		labels: labels,
	}

	globalMu.Lock()
	globalMetrics = append(globalMetrics, dg)
	globalMu.Unlock()

	return dg
}

func (g *DGauge) Set(value float64, values ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(values) != len(g.labels) {
		return
	}

	idx := -1
	for i, m := range g.metrics {
		if slices.Equal(values, m.LabelValues) {
			idx = i
			break
		}
	}

	if idx != -1 {
		g.metrics[idx].Value = value
	} else {
		g.metrics = append(g.metrics, DGaugeMetric{
			Value:       value,
			LabelValues: values,
		})
	}
}

func (g *DGauge) Reset(values ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	idx := -1
	for i, m := range g.metrics {
		if slices.Equal(values, m.LabelValues) {
			idx = i
			break
		}
	}

	if idx != -1 {
		g.metrics = slices.Delete(g.metrics, idx, idx)
	}
}

func (g *DGauge) Get() []metric {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]metric, 0, len(g.metrics))
	for _, m := range g.metrics {
		labels := make(map[string]string, len(g.labels))
		for i, name := range g.labels {
			labels[name] = m.LabelValues[i]
		}

		result = append(result, metric{
			Name:   g.name,
			Labels: labels,
			Type:   "DGAUGE",
			Value:  m.Value,
		})
	}

	return result
}

type CounterMetric struct {
	Value       int64
	LabelValues []string
}

type Counter struct {
	name   string
	labels []string

	metrics []CounterMetric
	mu      sync.RWMutex
}

func NewCounter(name string, labels ...string) *Counter {
	c := &Counter{
		name:   name,
		labels: labels,
	}

	globalMu.Lock()
	globalMetrics = append(globalMetrics, c)
	globalMu.Unlock()

	return c
}

func (g *Counter) Inc(values ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(values) != len(g.labels) {
		return
	}

	idx := -1
	for i, m := range g.metrics {
		if slices.Equal(values, m.LabelValues) {
			idx = i
			break
		}
	}

	if idx != -1 {
		g.metrics[idx].Value += 1
	} else {
		g.metrics = append(g.metrics, CounterMetric{
			Value:       1,
			LabelValues: values,
		})
	}
}

func (g *Counter) Reset(values ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	idx := -1
	for i, m := range g.metrics {
		if slices.Equal(values, m.LabelValues) {
			idx = i
			break
		}
	}

	if idx != -1 {
		g.metrics = slices.Delete(g.metrics, idx, idx)
	}
}

func (g *Counter) Get() []metric {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]metric, 0, len(g.metrics))
	for _, m := range g.metrics {
		labels := make(map[string]string, len(g.labels))
		for i, name := range g.labels {
			labels[name] = m.LabelValues[i]
		}

		result = append(result, metric{
			Name:   g.name,
			Labels: labels,
			Type:   "COUNTER",
			Value:  float64(m.Value),
		})
	}

	return result
}
