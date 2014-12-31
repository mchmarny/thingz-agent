package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
)

// CPUProvider is the provider for CPU information
type CPUProvider struct {
	Config *ProviderConfig
}

// Provide CPU metrics
func (p CPUProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		cpu := sigar.Cpu{}
		if err := cpu.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			return err
		}

		col := types.NewMetricCollection(p.Config.Source, p.Config.Dimension)

		col.Add(types.NewMetricSample("total", cpu.Total()))
		col.Add(types.NewMetricSample("user", cpu.User))
		col.Add(types.NewMetricSample("nice", cpu.Nice))
		col.Add(types.NewMetricSample("sys", cpu.Sys))
		col.Add(types.NewMetricSample("idle", cpu.Idle))
		col.Add(types.NewMetricSample("wait", cpu.Wait))

		out <- col

	}

	return nil

}
