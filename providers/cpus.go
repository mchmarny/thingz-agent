package providers

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz-commons/types"
)

const (
	INDEX_SEPERATOR = "_"
)

// CPUSProvider is the provider for CPU information
type CPUSProvider struct {
	Config *ProviderConfig
}

// Provide CPU metrics
func (p CPUSProvider) Provide(out chan<- *types.MetricCollection) error {

	ticker := time.NewTicker(p.Config.Frequency)

	for t := range ticker.C {

		cpul := sigar.CpuList{}
		if err := cpul.Get(); err != nil {
			log.Fatalf("Error in %v execution: %v", t, err)
			return err
		}

		col := types.NewMetricCollection(p.Config.Source, p.Config.Dimension)

		for i, c := range cpul.List {
			col.Add(types.NewMetricSample(fmt.Sprintf("total%s%d", INDEX_SEPERATOR, i), c.Total()))
			col.Add(types.NewMetricSample(fmt.Sprintf("user%s%d", INDEX_SEPERATOR, i), c.User))
			col.Add(types.NewMetricSample(fmt.Sprintf("nice%s%d", INDEX_SEPERATOR, i), c.Nice))
			col.Add(types.NewMetricSample(fmt.Sprintf("sys%s%d", INDEX_SEPERATOR, i), c.Sys))
			col.Add(types.NewMetricSample(fmt.Sprintf("idle%s%d", INDEX_SEPERATOR, i), c.Idle))
			col.Add(types.NewMetricSample(fmt.Sprintf("wait%s%d", INDEX_SEPERATOR, i), c.Wait))
		}

		out <- col

	}

	return nil

}
