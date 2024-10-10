package generic

import "github.com/Sinketsu/artifactsmmo/internal/monitoring"

var (
	goldCount      = monitoring.NewDGauge("gold_count", "character")
	tasksCoinCount = monitoring.NewDGauge("tasks_coin_count", "character")
	requestCount   = monitoring.NewCounter("request_count")
	skillLevel     = monitoring.NewDGauge("skill_level", "character", "skill")
)
