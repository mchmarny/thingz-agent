package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mchmarny/thingz/providers"
	"github.com/mchmarny/thingz/publishers"
	"github.com/mchmarny/thingz/types"
)

const (
	FORMAT_ERROR = "Invalid strategy format: "

	STRATEGY_CPU  = "cpu"
	STRATEGY_MEM  = "mem"
	STRATEGY_SWAP = "swap"
	STRATEGY_LOAD = "load"

	PUB_CONSOLE = "stdout"
)

func newCollector() (*collector, error) {

	c := &collector{
		providers: make(map[string]providers.Provider),
	}

	// TODO: Derive publisher dynamically
	if conf.Publisher == PUB_CONSOLE {
		c.publisher = publishers.ConsolePublisher{}
	}

	for i, p := range strings.Split(conf.Strategy, ",") {

		log.Printf("Parsing strategy[%d]:%s", i, p)
		strategy := strings.Split(strings.Trim(p, " "), ":")
		if len(strategy) != 2 {
			log.Fatal(FORMAT_ERROR)
			return nil, errors.New(FORMAT_ERROR + p)
		}

		// frequency of execution
		n, err := strconv.Atoi(strategy[1])
		if err != nil {
			log.Fatal(err)
			return nil, errors.New(FORMAT_ERROR + p)
		}

		freq := time.Duration(int32(n)) * time.Second
		group := strings.ToLower(strings.Trim(strategy[0], " "))

		log.Printf("Loading %s provider", group)

		switch group {
		case STRATEGY_CPU:
			c.providers[group] = providers.CPUProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_MEM:
			c.providers[group] = providers.MemoryProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_SWAP:
			c.providers[group] = providers.SwapProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_LOAD:
			c.providers[group] = providers.LoadProvider{
				Group:     group,
				Frequency: freq,
			}
		default:
			log.Fatal(FORMAT_ERROR)
			return nil, errors.New(FORMAT_ERROR + p)
		}

	}

	return c, nil
}

// collector type
type collector struct {
	providers map[string]providers.Provider
	publisher publishers.Publisher
}

//
func (c *collector) collect() error {

	chCount := 0

	log.Println("Providers:")
	for k, p := range c.providers {
		d, err := p.Describe()
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Printf("   %s - %v", k, d)
		chCount += len(d.Metrics)
	}

	log.Printf("Creating %d channels", chCount)

	errCh := make(chan error, 1)
	metricCh := make(chan *types.Metric, chCount)

	// start the collection routines
	for _, p := range c.providers {
		log.Printf("Starting collector: %v", p)
		go p.Provide(metricCh)
	}

	for {
		select {
		case err := <-errCh:
			log.Printf("Error: %v", err)
		case metric := <-metricCh:
			c.publisher.Publish(metric)
		default:
			// nothing to do
		}
	}

	return nil

}
