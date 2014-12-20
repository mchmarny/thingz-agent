package main

import (
	"log"

	"github.com/mchmarny/thingz/providers"
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

}
