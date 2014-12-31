package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons"
)

// LoadProvider is the provider for load information
type LoadProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p LoadProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Provide load metrics
func (p LoadProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Frequency)

	for t := range ticker.C {

		src := sigar.LoadAverage{}
		if err := src.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		col := types.NewMetricCollection(p.Group, t)

		col.Add(types.NewMetric(p.Group, "min1", src.One))
		col.Add(types.NewMetric(p.Group, "min5", src.Five))
		col.Add(types.NewMetric(p.Group, "min15", src.Fifteen))

		out <- col

	}

	return nil

}
