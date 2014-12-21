package providers

import (
	"log"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/mchmarny/thingz/types"
)

// LoadProvider is the provider for load information
type LoadProvider struct {
	Group     string
	Frequency time.Duration
}

// SetFrequency of execution
func (p LoadProvider) SetFrequency(f time.Duration) {
	p.Frequency = f
}

// Describe the load metric provider capabilities
func (p LoadProvider) Describe() (*types.Metadata, error) {

	m := types.NewMetadata(p.Group)

	// total CPU
	m.AddMetric("min1", "Average load for 1 min")
	m.AddMetric("min5", "Average load for 5 min")
	m.AddMetric("min15", "Average load for 15 min")

	return m, nil
}

// Provide load metrics
func (p LoadProvider) Provide(out chan<- *types.Metric) error {

	ticker := time.NewTicker(p.Frequency)

	for t := range ticker.C {

		log.Println("%s provider at %v", p.Group, t)

		src := sigar.LoadAverage{}
		if err := src.Get(); err != nil {
			log.Fatal(err)
			return err
		}

		go func() {
			out <- types.NewMetric(p.Group, "min1", src.One)
		}()
		go func() {
			out <- types.NewMetric(p.Group, "min5", src.Five)
		}()
		go func() {
			out <- types.NewMetric(p.Group, "min15", src.Fifteen)
		}()

	}

	return nil

}
