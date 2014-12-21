package types

import "fmt"

// NewMetadata is a factory method for Metadata
func NewMetadata(group string) *Metadata {
	return &Metadata{
		Group:   group,
		Metrics: make(map[string]string),
	}
}

// Metadata represents a metric description
type Metadata struct {

	// Group name
	Group string `json:"group"`

	// Context data for this metric
	Metrics map[string]string `json:"metrics"`
}

// AddMetric adds metric detail to metadata
// and return itself to allow for chaining
func (m *Metadata) AddMetric(name, desc string) *Metadata {
	m.Metrics[name] = desc
	return m
}

func (m *Metadata) String() string {
	return fmt.Sprintf("Metadata: %s [%v]", m.Group, m.Metrics)
}
