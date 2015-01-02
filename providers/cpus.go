package providers

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons"
)

const (
	INDEX_SEPERATOR = "_"
)

// CPUSProvider is the provider for CPU information
type CPUSProvider struct {
	Config *ProviderConfig
}

// Provide CPU metrics
func (p CPUSProvider) Provide(out chan<- *commons.Metric, outErr chan<- error) {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		cpul := sigar.CpuList{}
		if err := cpul.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			outErr <- err
		}

		for i, c := range cpul.List {
			out <- commons.NewMetric(
				p.Config.Source,
				p.Config.Dimension,
				fmt.Sprintf("total%s%d", INDEX_SEPERATOR, i),
				c.Total(),
			)

			out <- commons.NewMetric(
				p.Config.Source,
				p.Config.Dimension,
				fmt.Sprintf("user%s%d", INDEX_SEPERATOR, i),
				c.User,
			)

			out <- commons.NewMetric(
				p.Config.Source,
				p.Config.Dimension,
				fmt.Sprintf("nice%s%d", INDEX_SEPERATOR, i),
				c.Nice,
			)

			out <- commons.NewMetric(
				p.Config.Source,
				p.Config.Dimension,
				fmt.Sprintf("sys%s%d", INDEX_SEPERATOR, i),
				c.Sys,
			)

			out <- commons.NewMetric(
				p.Config.Source,
				p.Config.Dimension,
				fmt.Sprintf("idle%s%d", INDEX_SEPERATOR, i),
				c.Idle,
			)

			out <- commons.NewMetric(
				p.Config.Source,
				p.Config.Dimension,
				fmt.Sprintf("wait%s%d", INDEX_SEPERATOR, i),
				c.Wait,
			)
		}

	}

}
