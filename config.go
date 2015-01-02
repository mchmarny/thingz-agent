package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mchmarny/thingz-agent/publishers"
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	log.Printf("Initializing thingz-agent on %s...", hostname)

	flag.StringVar(&conf.Strategy, "strategy", "cpu:1,mem:1,swap:5,load:5", "Provider strategy")
	flag.StringVar(&conf.Source, "source", hostname, "Event source")
	flag.StringVar(&conf.Publisher, "publisher", publishers.PUB_CONSOLE, "Publishing target")
	flag.StringVar(&conf.PublisherArgs, "publisher-args", "", "Publishing arguments")
	flag.BoolVar(&conf.Verbose, "verbose", false, "Verbose outpur")

	flag.Parse()

	if strings.Index(conf.Source, ".") > -1 {
		conf.Source = strings.Replace(conf.Source, ".", "-", -1)
		log.Printf(
			"Source included offending characters [.] using [%s] instead",
			conf.Source,
		)
	}

}

// conf is a global instance of Config
var conf = &Config{}

// Config
type Config struct {
	Source        string
	Strategy      string
	Publisher     string
	PublisherArgs string
	Verbose       bool
}
