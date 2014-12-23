package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mchmarny/thingz-agent/providers"
	"github.com/mchmarny/thingz-agent/publishers"
	"github.com/mchmarny/thingz-agent/types"
)

const (
	STRATEGY_CPU  = "cpu"
	STRATEGY_CPUS = "cpus"
	STRATEGY_MEM  = "mem"
	STRATEGY_SWAP = "swap"
	STRATEGY_LOAD = "load"
	PUB_CONSOLE   = "stdout"
)

func newCollector() (*collector, error) {

	c := &collector{
		providers: make(map[string]providers.Provider),
	}

	var pub publishers.Publisher
	var err error
	if conf.Publisher == PUB_CONSOLE {
		log.Println("Using console publisher")
		pub, err = publishers.NewConsolePublisher()
	} else {
		log.Println("Using InfluxDB publisher")
		pub, err = publishers.NewInfluxDBPublisher(
			conf.Source,
			conf.Publisher,
		)
	}
	if err != nil {
		log.Fatalln("Error while creating publisher")
		return nil, err
	}
	c.publisher = pub

	for _, p := range strings.Split(conf.Strategy, ",") {

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

		switch group {
		case STRATEGY_CPU:
			c.providers[group] = providers.CPUProvider{
				Group:     group,
				Frequency: freq,
			}
		case STRATEGY_CPUS:
			c.providers[group] = providers.CPUSProvider{
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

	errCh := make(chan error, 1)
	metricCh := make(chan *types.MetricCollection, len(c.providers))

	// start the collection routines
	for _, p := range c.providers {
		log.Printf("Starting collector for %v", p)
		go p.Provide(metricCh)
	}

	for {
		select {
		case err := <-errCh:
			log.Printf("Error: %v", err)
		case col := <-metricCh:
			go c.publisher.Publish(col)
		default:
			// nothing to do
		}
	}

	return nil

}
