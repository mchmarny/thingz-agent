package providers

import (
	"github.com/mchmarny/thingz/types"
)

// Provider describes the metric provider functionality
type Provider interface {

	// Describe provider capabilities
	Describe() (*types.Metadata, error)

	// Provide metric group
	Provide(freq int, out <-chan *types.Metric) error
}
