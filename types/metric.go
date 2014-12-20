package types

import (
	"time"
)

// NewMetric is a factory method for Metric
func NewMetric(dimension string, value interface{}) *Metric {
	return &Metric{
		Timestamp: time.Now(),
		Dimension: dimension,
		Value:     value,
	}
}

// Metric represents a generic metric collection event
type Metric struct {

	// Timestamp of when the metric was captured
	Timestamp time.Time `json:"ts"`

	// Dimension this metric represents
	Dimension string `json:"dim"`

	// Value of this metric
	Value interface{} `json:"val"`

	// Unit of this metric
	Unit string `json:"unit,omitempty"`

	// Context data for this metric
	Context map[string]string `json:"ctx,omitempty"`
}

// SetUnit sets metric unit
// and return itself to allow for chaining
func (m *Metric) SetUnit(unit string) *Metric {
	m.Unit = unit
	return m
}

// AddContext adds context to this metric
// and return itself to allow for chaining
func (m *Metric) AddContext(key, val string) *Metric {

	if m.Context == nil {
		m.Context = make(map[string]string)
	}

	m.Context[key] = val
	return m
}
