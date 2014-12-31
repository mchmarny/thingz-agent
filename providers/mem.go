package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
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
