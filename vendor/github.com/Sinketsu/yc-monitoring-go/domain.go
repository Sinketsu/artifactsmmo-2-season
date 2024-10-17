package ycmonitoringgo

const (
	TYPE_DGAUGE  string = "DGAUGE"
	TYPE_IGAUGE  string = "IGAUGE"
	TYPE_COUNTER string = "COUNTER"
	TYPE_RATE    string = "RATE"
)

type Request struct {
	Metrics []metric `json:"metrics"`
}

type Metric interface {
	Name() string
	GetMetrics() []metric
}

type metric struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
	Type   string            `json:"type"`
	Value  any               `json:"value"`
}
