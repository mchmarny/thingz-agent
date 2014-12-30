package types

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// NewMetric is a factory method for MetricCollection
func NewMetricCollection(group string, at time.Time) *MetricCollection {
	var list []*Metric
	return &MetricCollection{
		Runtime: at,
		Group:   group,
		Metrics: list,
	}
}

// MetricCollection represents a generic metric collection event
type MetricCollection struct {

	// Runtime of when the metric was captured
	Runtime time.Time `json:"t"`

	// Group this metric represents
	Group string `json:"group"`

	// Metrics represents a collection of metric items
	Metrics []*Metric `json:"metrics"`
}

// Add adds metric to collection
// and return itself to allow for chaining
func (m *MetricCollection) Add(item *Metric) {
	m.Metrics = append(m.Metrics, item)
}

func (m *MetricCollection) String() string {
	return fmt.Sprintf(
		"MetricCollection: [ Group:%s, Runtime:%s, Metrics:%v ]",
		m.Group, m.Runtime, m.Metrics,
	)
}

// ToBytes converts content of the current message into byte array
func (m *MetricCollection) ToBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("unable to marshal: %v", err.Error())
	}
	return b
}
