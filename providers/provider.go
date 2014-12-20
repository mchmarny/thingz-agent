package providers

import (
	"github.com/mchmarny/thingz/types"
	"time"
)

// Provider describes the metric provider functionality
type Provider interface {

	// Describe provider capabilities
	Describe() (*types.Metadata, error)

	// Provide metric group
	Provide(freq time.Duration, out <-chan *types.Metric) error
}
