package providers

import (
	"fmt"
	"log"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz/types"
)

const (
	GROUP_NAME = "CPU"
)

// CPUProvider is the provider for CPU information
type CPUProvider struct{}

// Get CPU metric group
func (p *CPUProvider) Get() (*types.MetricGroup, error) {

	g := types.NewMetricGroup(GROUP_NAME)

	cpu := sigar.Cpu{}
	if err := cpu.Get(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	g.Add(types.NewMetric("Total", cpu.Total()))
	g.Add(types.NewMetric("User", cpu.User))
	g.Add(types.NewMetric("Nice", cpu.Nice))
	g.Add(types.NewMetric("Sys", cpu.Sys))
	g.Add(types.NewMetric("Idle", cpu.Idle))
	g.Add(types.NewMetric("Wait", cpu.Wait))

	cpul := sigar.CpuList{}
	if err := cpul.Get(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	for i, c := range cpul.List {
		g.Add(types.NewMetric(fmt.Sprintf("C%d-Total", i), c.Total()))
		g.Add(types.NewMetric(fmt.Sprintf("C%d-User", i), c.User))
		g.Add(types.NewMetric(fmt.Sprintf("C%d-Nice", i), c.Nice))
		g.Add(types.NewMetric(fmt.Sprintf("C%d-Sys", i), c.Sys))
		g.Add(types.NewMetric(fmt.Sprintf("C%d-Idle", i), c.Idle))
		g.Add(types.NewMetric(fmt.Sprintf("C%d-Wait", i), c.Wait))
	}

	return g, nil
}
