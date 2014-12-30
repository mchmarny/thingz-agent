package publishers

import (
	"github.com/mchmarny/thingz-agent/types"
)

// Publisher describes the metric publisher functionality
type Publisher interface {

	// Publish metric
	Publish(m *types.MetricCollection)

	// Finalize tells the publisher to close used resources
	// and do any general cleanup it needs
	Finalize()
}
