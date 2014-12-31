package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons"
)

// LoadProvider is the provider for load information
type LoadProvider struct {
	Config *ProviderConfig
}

// Provide load metrics
func (p LoadProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		src := sigar.LoadAverage{}
		if err := src.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			return err
		}

		col := types.NewMetricCollection(p.Config.Source, p.Config.Dimension)

		col.Add(types.NewMetricSample("min1", src.One))
		col.Add(types.NewMetricSample("min5", src.Five))
		col.Add(types.NewMetricSample("min15", src.Fifteen))

		out <- col

	}

	return nil

}
