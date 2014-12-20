package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mchmarny/thingz/providers"
	"github.com/mchmarny/thingz/types"
)

func main() {
	conf.printHeader()

	c := &providers.CPUProvider{}

	d, err := c.Describe()

	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Provider: %s", d.Group)

	for k, n := range d.Metrics {
		log.Printf("   %s - %s", k, n)
	}

	// make sure we can shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	doneDone := make(chan bool)
	signal.Notify(sigCh, os.Interrupt)
	errCh := make(chan error, 1)

	// worker channels
	provCh := make(chan *types.Metric, len(d.Metrics))

	go func() {
		freq := time.Second * 1
		errCh <- c.Provide(freq, provCh)
	}()

	go func() {

		for {
			select {
			case err := <-errCh:
				log.Printf("Error: %v", err)
			case res := <-provCh:
				log.Printf("Result: %v", res)
			case ex := <-sigCh:
				log.Println("Shutting down due to: ", ex)
				doneDone <- true
			default:
				// nothing to do
			}
		}

	}()
	<-doneDone

}
