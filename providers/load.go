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
func (p LoadProvider) Provide(out chan<- *commons.Metric, outErr chan<- error) {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		src := sigar.LoadAverage{}
		if err := src.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			outErr <- err
		}

		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "min1", src.One)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "min5", src.Five)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "min15", src.Fifteen)

	}

}
