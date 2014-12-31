package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mchmarny/thingz-agent/publishers"
)

const (
	APP_VERSION = "v0.3"
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&conf.Strategy, "strategy", "cpu:1,mem:1,swap:5,load:5", "Provider strategy")
	flag.StringVar(&conf.Source, "source", hostname, "Event source")
	flag.StringVar(&conf.Publisher, "publisher", publishers.PUB_CONSOLE, "Publishing target")
	flag.StringVar(&conf.PublisherArgs, "publisher-args", "", "Publishing arguments")
	flag.BoolVar(&conf.Verbose, "verbose", false, "Verbose outpur")

	flag.Parse()

	conf.Version = APP_VERSION
	log.Printf("Version: %s", APP_VERSION)

	if strings.Index(conf.Source, ".") > -1 {
		conf.Source = strings.Replace(conf.Source, ".", "-", -1)
		log.Printf(
			"Source included offending characters [.] using [%s] instead",
			conf.Source,
		)
	}

}

var conf = &Config{}

type Config struct {
	Version       string
	Source        string
	Strategy      string
	Publisher     string
	PublisherArgs string
	Verbose       bool
}
