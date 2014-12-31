package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
)

// CPUProvider is the provider for CPU information
type CPUProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p CPUProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Provide CPU metrics
func (p CPUProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Frequency)

	for t := range ticker.C {

		cpu := sigar.Cpu{}
		if err := cpu.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		col := types.NewMetricCollection(p.Group, t)

		col.Add(types.NewMetric(p.Group, "total", cpu.Total()))
		col.Add(types.NewMetric(p.Group, "user", cpu.User))
		col.Add(types.NewMetric(p.Group, "nice", cpu.Nice))
		col.Add(types.NewMetric(p.Group, "sys", cpu.Sys))
		col.Add(types.NewMetric(p.Group, "idle", cpu.Idle))
		col.Add(types.NewMetric(p.Group, "wait", cpu.Wait))

		out <- col

	}

	return nil

}
