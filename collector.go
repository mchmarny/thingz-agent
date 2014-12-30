package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/mchmarny/thingz-agent/providers"
	"github.com/mchmarny/thingz-agent/publishers"
	"github.com/mchmarny/thingz-agent/types"
)

const (
	// TODO: refactor for each provider to describe itself
	STRATEGY_CPU  = "cpu"
	STRATEGY_CPUS = "cpus"
	STRATEGY_MEM  = "mem"
	STRATEGY_SWAP = "swap"
	STRATEGY_LOAD = "load"
	STRATEGY_PROC = "proc"

	PUB_CONSOLE  = "stdout"
	PUB_INFLUXDB = "influxdb"
	PUB_KAFKA    = "kafka"
)

func newCollector() (*collector, error) {

	// Load publisher
	var pub publishers.Publisher
	var err error

	switch conf.Publisher {
	case PUB_CONSOLE:
		pub, err = publishers.NewConsolePublisher()
	case PUB_INFLUXDB:
		pub, err = publishers.NewInfluxDBPublisher(
			conf.Source,
			conf.PublisherArgs,
			conf.Verbose,
		)
	case PUB_KAFKA:
		pub, err = publishers.NewKafkaPublisher(
			conf.Source,
			conf.PublisherArgs,
		)
	default:
		errors.New("Invalid publishing target: " + conf.Publisher)
	}

	if err != nil {
		log.Fatalln("Error while creating publisher")
		return nil, err
	}

	// Load collector
	c := &collector{
		providers: make(map[string]providers.Provider),
		publisher: pub,
	}

	for _, p := range strings.Split(conf.Strategy, ",") {

		// get strategy
		strategy := strings.Split(strings.Trim(p, " "), ":")
		if len(strategy) != 2 {
			log.Fatal(FORMAT_ERROR)
			return nil, errors.New(FORMAT_ERROR + p)
		}

		// frequency of execution for this strategy
		n, err := strconv.Atoi(strategy[1])
		if err != nil {
			log.Fatal(err)
			return nil, errors.New(FORMAT_ERROR + p)
		}

		freq := time.Duration(int32(n)) * time.Second
		group := strings.ToLower(strings.Trim(strategy[0], " "))

		// TODO: spool these into a map first
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
		case STRATEGY_PROC:
			c.providers[group] = providers.ProcProvider{
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

// collect
func (c *collector) collect() error {

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	metricCh := make(chan *types.MetricCollection, len(c.providers))

	signal.Notify(sigCh, os.Interrupt)

	// start the collection routines
	for _, p := range c.providers {
		go p.Provide(metricCh)
	}

	// watch magic happen
	for {
		select {
		case err := <-errCh:
			log.Printf("Collector error: %v", err)
		case col := <-metricCh:
			// publish collection upon receiving
			go c.publisher.Publish(col)
		case sig := <-sigCh:
			log.Printf("Notifying publisher: %v", sig)
			c.publisher.Finalize()
		default:
			// nothing to do
		}
	}

	return nil

}
