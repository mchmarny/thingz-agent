package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons"
)

// MemoryProvider is the provider for memory information
type MemoryProvider struct {
	Config *ProviderConfig
}

// Provide memory metrics
func (p MemoryProvider) Provide(out chan<- *commons.Metric, outErr chan<- error) {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		src := sigar.Mem{}
		if err := src.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			outErr <- err
		}

		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "free", src.Free)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "used", src.Used)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "actual-free", src.ActualFree)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "actual-used", src.ActualUsed)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "total", src.Total)

	}

}
