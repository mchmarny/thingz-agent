package types

import (
	"fmt"
	"time"
)

// NewMetric is a factory method for Metric
func NewMetric(group, dimension string, value interface{}) *Metric {
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

func (m *Metric) String() string {
	return fmt.Sprintf(
		"Metric: [ Dimension:%s, Timestamp:%v, Value:%v, Unit:%s, Context:%v ]",
		m.Dimension, m.Timestamp, m.Value, m.Unit, m.Context,
	)
}
