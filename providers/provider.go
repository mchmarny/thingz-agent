package providers

import (
	"github.com/mchmarny/thingz/types"
)

// Provider describes the metric provider functionality
type Provider interface {

	// Get metric group
	Get() (types.MetricGroup, error)
}
