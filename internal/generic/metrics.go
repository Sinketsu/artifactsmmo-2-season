package generic

import ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"

var (
	goldCount      = ycmonitoringgo.NewDGauge("gold_count", ycmonitoringgo.DefaultRegistry, "character")
	tasksCoinCount = ycmonitoringgo.NewDGauge("tasks_coin_count", ycmonitoringgo.DefaultRegistry, "character")
	requestCount   = ycmonitoringgo.NewCounter("request_count", ycmonitoringgo.DefaultRegistry)
	skillLevel     = ycmonitoringgo.NewDGauge("skill_level", ycmonitoringgo.DefaultRegistry, "character", "skill")
)
