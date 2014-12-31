package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
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
