package providers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-agent/types"
)

// ProcProvider is the provider for CPU information
type ProcProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p ProcProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Provide CPU metrics
func (p ProcProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Frequency)

	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
		return err
	}

	for t := range ticker.C {

		pids := sigar.ProcList{}
		if err := pids.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		col := types.NewMetricCollection(p.Group, t)

		for _, pid := range pids.List {
			state := sigar.ProcState{}
			mem := sigar.ProcMem{}

			if err := state.Get(pid); err != nil {
				continue
			}
			if err := mem.Get(pid); err != nil {
				continue
			}

			safe := reg.ReplaceAllString(state.Name, "-")
			safe = strings.ToLower(strings.Trim(safe, "-"))

			col.Add(
				types.NewMetric(p.Group,
					fmt.Sprintf("p%d-%s", pid, safe),
					mem.Resident/1024),
			)

		}

		out <- col

	}

	return nil

}
