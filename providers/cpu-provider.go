package providers

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz/types"
)

const (
	GROUP_NAME = "CPU"
)

// CPUProvider is the provider for CPU information
type CPUProvider struct{}

// Describe the CPU metric provider capabilities
func (p CPUProvider) Describe() (*types.Metadata, error) {

	m := types.NewMetadata(GROUP_NAME)

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
func (p CPUProvider) Provide(freq time.Duration, out chan<- *types.Metric) error {

	ticker := time.NewTicker(freq)

	for t := range ticker.C {

		log.Println("Provider at: ", t)

		cpu := sigar.Cpu{}
		if err := cpu.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		out <- types.NewMetric("total", cpu.Total())
		out <- types.NewMetric("user", cpu.User)
		out <- types.NewMetric("nice", cpu.Nice)
		out <- types.NewMetric("sys", cpu.Sys)
		out <- types.NewMetric("idle", cpu.Idle)
		out <- types.NewMetric("wait", cpu.Wait)

		cpul := sigar.CpuList{}
		if err := cpul.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		for i, c := range cpul.List {
			out <- types.NewMetric(fmt.Sprintf("c%d-total", i), c.Total())
			out <- types.NewMetric(fmt.Sprintf("c%d-user", i), c.User)
			out <- types.NewMetric(fmt.Sprintf("c%d-nice", i), c.Nice)
			out <- types.NewMetric(fmt.Sprintf("c%d-sys", i), c.Sys)
			out <- types.NewMetric(fmt.Sprintf("c%d-idle", i), c.Idle)
			out <- types.NewMetric(fmt.Sprintf("c%d-wait", i), c.Wait)
		}

	}

	return nil

}
