package providers

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-agent/types"
)

const (
	INDEX_SEPERATOR = "_"
)

// CPUSProvider is the provider for CPU information
type CPUSProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p CPUSProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Describe the CPU metric provider capabilities
func (p CPUSProvider) Describe() (*types.Metadata, error) {

	l := sigar.CpuList{}
	if err := l.Get(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	m := types.NewMetadata(p.Group)

	// individual CPUs
	for i, _ := range l.List {

		m.AddMetric(fmt.Sprintf("total%s%d", INDEX_SEPERATOR, i), fmt.Sprintf("Total combined for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("user%s%d", INDEX_SEPERATOR, i), fmt.Sprintf("User time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("nice%s%d", INDEX_SEPERATOR, i), fmt.Sprintf("Nice time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("sys%s%d", INDEX_SEPERATOR, i), fmt.Sprintf("Sys time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("idle%s%d", INDEX_SEPERATOR, i), fmt.Sprintf("Idle time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("wait%s%d", INDEX_SEPERATOR, i), fmt.Sprintf("Wait time for CPU[%d]", i))
	}

	return m, nil
}

// Provide CPU metrics
func (p CPUSProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Frequency)

	for t := range ticker.C {

		cpul := sigar.CpuList{}
		if err := cpul.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		col := types.NewMetricCollection(p.Group, t)

		for i, c := range cpul.List {
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("total%s%d", INDEX_SEPERATOR, i), c.Total()))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("user%s%d", INDEX_SEPERATOR, i), c.User))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("nice%s%d", INDEX_SEPERATOR, i), c.Nice))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("sys%s%d", INDEX_SEPERATOR, i), c.Sys))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("idle%s%d", INDEX_SEPERATOR, i), c.Idle))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("wait%s%d", INDEX_SEPERATOR, i), c.Wait))
		}

		out <- col

	}

	return nil

}
