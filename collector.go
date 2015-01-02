package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/mchmarny/thingz-agent/providers"
	"github.com/mchmarny/thingz-agent/publishers"
	"github.com/mchmarny/thingz-commons"
)

func newCollector() (*collector, error) {

	// Load publisher
	pub, err := publishers.GetPublisher(
		conf.Source,
		conf.Publisher,
		conf.PublisherArgs,
	)
	if err != nil {
		log.Fatalln("Error while creating publisher")
		return nil, err
	}

	provs, err := providers.GetProviders(conf.Source, conf.Strategy)
	if err != nil {
		log.Fatalln("Error while creating providers")
		return nil, err
	}

	// Create a collector
	c := &collector{
		publisher: pub,
		providers: provs,
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

	sigCh := make(chan os.Signal, 1)
	metricCh := make(chan *commons.Metric, len(c.providers))
	prvErrCh := make(chan error, len(c.providers))
	pubErrCh := make(chan error, len(c.providers))

	signal.Notify(sigCh, os.Interrupt)

	// start the collection routines
	for _, p := range c.providers {
		go p.Provide(metricCh, prvErrCh)
	}

	// send provider output to publisher
	go c.publisher.Publish(metricCh, pubErrCh)

	// watch magic happen... or shit fall apart
	for {
		select {
		case prvErr := <-prvErrCh:
			log.Printf("Provider error: %v", prvErr)
		case pubErr := <-pubErrCh:
			log.Printf("Publisher error: %v", pubErr)
		case sig := <-sigCh:
			log.Printf("Finalizing %v due to %v", c.publisher, sig)
			c.publisher.Finalize()
		default:
			// nothing to do
		}
	}

	return nil

}
