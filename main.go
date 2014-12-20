package main

import (
	"log"

	"github.com/mchmarny/thingz/providers"
)

func main() {
	conf.printHeader()

	p := &providers.CPUProvider{}
	g, err := p.Get()

	if err != nil {
		log.Fatal(err)
	}

	log.Print(g)

}
