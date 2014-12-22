package main

import (
	"flag"
	"log"
	"os"
)

const (
	APP_VERSION  = "0.0.1"
	FORMAT_ERROR = "Invalid strategy format: "
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	hostname, _ := os.Hostname()

	flag.StringVar(&conf.Strategy, "strategy", "cpu:1,mem:1,swap:5,load:5", "Provider strategy")
	flag.StringVar(&conf.Source, "source", hostname, "Event source")
	flag.StringVar(&conf.Publisher, "publisher", "stdout", "Publishing target")

	conf.Version = APP_VERSION

	flag.Parse()

}

var conf = &Config{}

type Config struct {
	Version   string
	Source    string
	Strategy  string
	Publisher string
}
