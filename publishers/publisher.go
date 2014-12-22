package publishers

import (
	"github.com/mchmarny/thingz-agent/types"
)

// Publisher describes the metric publisher functionality
type Publisher interface {

	// Publish metric
	Publish(m *types.MetricCollection) error
}
