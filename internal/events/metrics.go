package events

import ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"

var (
	errorRate  = ycmonitoringgo.NewRate("events_errors", ycmonitoringgo.DefaultRegistry)
	eventCount = ycmonitoringgo.NewIGauge("events_count", ycmonitoringgo.DefaultRegistry)
)
