package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-agent/types"
)

// MemoryProvider is the provider for memory information
type MemoryProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p MemoryProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Describe the memory metric provider capabilities
func (p MemoryProvider) Describe() (*types.Metadata, error) {

	m := types.NewMetadata(p.Group)

	// total CPU
	m.AddMetric("free", "Amount of free memory")
	m.AddMetric("used", "Amount of used memory")
	m.AddMetric("afree", "Amount of actual free memory")
	m.AddMetric("aused", "Amount of actual used memory")
	m.AddMetric("total", "Amount of total memory")

	return m, nil
}

// Provide memory metrics
func (p MemoryProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Frequency)

	for t := range ticker.C {

		src := sigar.Mem{}
		if err := src.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		col := types.NewMetricCollection(p.Group, t)

		col.Add(types.NewMetric(p.Group, "free", src.Free))
		col.Add(types.NewMetric(p.Group, "used", src.Used))
		col.Add(types.NewMetric(p.Group, "actual-free", src.ActualFree))
		col.Add(types.NewMetric(p.Group, "actual-used", src.ActualUsed))
		col.Add(types.NewMetric(p.Group, "total", src.Total))

		out <- col

	}

	return nil

}
