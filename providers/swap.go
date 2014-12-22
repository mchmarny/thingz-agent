package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz/types"
)

// SwapProvider is the provider for swap information
type SwapProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p SwapProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Describe the swap metric provider capabilities
func (p SwapProvider) Describe() (*types.Metadata, error) {

	m := types.NewMetadata(p.Group)

	// total CPU
	m.AddMetric("free", "Amount of free swap memory")
	m.AddMetric("used", "Amount of used swap memory")
	m.AddMetric("total", "Amount of total swap memory")

	return m, nil
}

// Provide swap metrics
func (p SwapProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Frequency)

	for t := range ticker.C {

		src := sigar.Swap{}
		if err := src.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		col := types.NewMetricCollection(p.Group, t)

		col.Add(types.NewMetric(p.Group, "free", src.Free))
		col.Add(types.NewMetric(p.Group, "used", src.Used))
		col.Add(types.NewMetric(p.Group, "total", src.Total))

		out <- col

	}

	return nil

}
