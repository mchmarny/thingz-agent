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
	Timestamp time.Time `json:"t"`

	// Dimension this metric represents
	Dimension string `json:"d"`

	// Value of this metric
	Value interface{} `json:"v"`

	// Unit of this metric
	Unit string `json:"u,omitempty"`
}

// SetUnit sets metric unit and return itself to allow for chaining
func (m *Metric) SetUnit(unit string) *Metric {
	m.Unit = unit
	return m
}
