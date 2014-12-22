package providers

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz/types"
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

// Describe the CPU metric provider capabilities
func (p CPUProvider) Describe() (*types.Metadata, error) {

	m := types.NewMetadata(p.Group)

	// total CPU
	m.AddMetric("total", "Total combined CPU")
	m.AddMetric("user", "User time")
	m.AddMetric("nice", "Nice time")
	m.AddMetric("sys", "Sys time")
	m.AddMetric("idle", "Idle time")
	m.AddMetric("wait", "Wait time")

	l := sigar.CpuList{}
	if err := l.Get(); err != nil {
		log.Fatal(err)
		return m, err
	}

	// individual CPUs
	for i, _ := range l.List {

		m.AddMetric(fmt.Sprintf("c%d-total", i), fmt.Sprintf("Total combined for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("c%d-user", i), fmt.Sprintf("User time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("c%d-nice", i), fmt.Sprintf("Nice time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("c%d-sys", i), fmt.Sprintf("Sys time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("c%d-idle", i), fmt.Sprintf("Idle time for CPU[%d]", i))
		m.AddMetric(fmt.Sprintf("c%d-wait", i), fmt.Sprintf("Wait time for CPU[%d]", i))
	}

	return m, nil
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

		cpul := sigar.CpuList{}
		if err := cpul.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		for i, c := range cpul.List {
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("c%d-total", i), c.Total()))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("c%d-user", i), c.User))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("c%d-nice", i), c.Nice))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("c%d-sys", i), c.Sys))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("c%d-idle", i), c.Idle))
			col.Add(types.NewMetric(p.Group, fmt.Sprintf("c%d-wait", i), c.Wait))
		}

		out <- col

	}

	return nil

}
