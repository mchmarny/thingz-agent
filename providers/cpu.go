package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons"
)

// CPUProvider is the provider for CPU information
type CPUProvider struct {
	Config *ProviderConfig
}

// Provide CPU metrics
func (p CPUProvider) Provide(out chan<- *commons.Metric, outErr chan<- error) {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		cpu := sigar.Cpu{}
		if err := cpu.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			outErr <- err
		}

		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "total", cpu.Total())
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "user", cpu.User)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "nice", cpu.Nice)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "sys", cpu.Sys)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "idle", cpu.Idle)
		out <- commons.NewMetric(p.Config.Source, p.Config.Dimension, "wait", cpu.Wait)

	}

}
