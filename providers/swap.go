package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
)

// SwapProvider is the provider for swap information
type SwapProvider struct {
	Config *ProviderConfig
}

// Provide swap metrics
func (p SwapProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		src := sigar.Swap{}
		if err := src.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			return err
		}

		col := types.NewMetricCollection(p.Config.Source, p.Config.Dimension)

		col.Add(types.NewMetricSample("free", src.Free))
		col.Add(types.NewMetricSample("used", src.Used))
		col.Add(types.NewMetricSample("total", src.Total))

		out <- col

	}

	return nil

}
