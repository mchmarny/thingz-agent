package providers

import (
	"github.com/mchmarny/thingz/types"
	"time"
)

// Provider describes the metric provider functionality
type Provider interface {

	// SetFrequency of execution
	SetFrequency(time.Duration)

	// Describe provider capabilities
	Describe() (*types.Metadata, error)

	// Provide metric group
	Provide(out chan<- *types.Metric) error
}
