package types

// NewMetricGroup is a factory method for MetricGroup
func NewMetricGroup(name string) *MetricGroup {
	return &MetricGroup{
		GroupName: name,
		Metrics:   []*Metric{},
	}
}

// MetricGroup represents a group of generic metric collection events
type MetricGroup struct {

	// GroupName of this groups
	GroupName string `json:"group"`

	// Value of this metric
	Metrics []*Metric `json:"list"`
}

// Add appends provided metric to the group
func (g *MetricGroup) Add(metric *Metric) {
	g.Metrics = append(g.Metrics, metric)
}
