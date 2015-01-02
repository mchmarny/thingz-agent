package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons"
)

// SwapProvider is the provider for swap information
type SwapProvider struct {
	Config *ProviderConfig
}

// Provide swap metrics
func (p SwapProvider) Provide(out chan<- *commons.Metric, outErr chan<- error) {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		src := sigar.Swap{}
		if err := src.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			outErr <- err
		}

		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "free", src.Free)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "used", src.Used)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "total", src.Total)

	}
}
