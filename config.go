package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	APP_VERSION = "0.0.1"
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	hostname, _ := os.Hostname()

	flag.BoolVar(&conf.Verbose, "verbose", false, "Debug info")
	flag.StringVar(&conf.Source, "source", hostname, "Event source")
	flag.StringVar(&conf.Strategy, "strategy", "cpu:1,mem:3", "Provider strategy")
	flag.StringVar(&conf.Publisher, "publisher", "stdout", "Publishing target")

	conf.Version = APP_VERSION

	flag.Parse()

}

var conf = &Config{}

type Config struct {
	Version   string
	Verbose   bool
	Source    string
	Strategy  string
	Publisher string
}

func (c *Config) printHeader() {
	fmt.Printf("Agent %s\n", conf.Version)
	fmt.Printf("   source: %s\n", conf.Source)
	fmt.Printf("   strategy: %s\n", conf.Strategy)
	fmt.Println("")
}
