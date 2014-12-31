package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
)

// MemoryProvider is the provider for memory information
type MemoryProvider struct {
	Config *ProviderConfig
}

// Provide memory metrics
func (p MemoryProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		src := sigar.Mem{}
		if err := src.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			return err
		}

		col := types.NewMetricCollection(p.Config.Source, p.Config.Dimension)

		col.Add(types.NewMetricSample("free", src.Free))
		col.Add(types.NewMetricSample("used", src.Used))
		col.Add(types.NewMetricSample("actual-free", src.ActualFree))
		col.Add(types.NewMetricSample("actual-used", src.ActualUsed))
		col.Add(types.NewMetricSample("total", src.Total))

		out <- col

	}

	return nil

}
