package providers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/cloudfoundry/gosigar"
	types "github.com/mchmarny/thingz-commons"
)

// ProcProvider is the provider for CPU information
type ProcProvider struct {
	Config *ProviderConfig
}

// Provide CPU metrics
func (p ProcProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Config.Frequency)

	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatalf("Error while creating regex: %v", err)
		return err
	}

	for t := range ticker.C {

		pids := sigar.ProcList{}
		if err := pids.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			return err
		}

		col := types.NewMetricCollection(p.Config.Source, p.Config.Dimension)

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
				types.NewMetricSample(
					fmt.Sprintf("p%d-%s", pid, safe),
					mem.Resident/1024),
			)

		}

		out <- col

	}

	return nil

}
