package generic

import "github.com/Sinketsu/artifactsmmo/internal/monitoring"

var (
	goldCount = monitoring.NewDGauge("gold_count", "character")
)
