package publishers

import (
	"log"
	"net/url"
	"strings"

	flux "github.com/influxdb/influxdb/client"
	"github.com/mchmarny/thingz-commons"
)

// NewInfluxDBPublisher parses connection string to InfluxDB
// and returned a configured version of the publisher
func NewInfluxDBPublisher(connStr string) (Publisher, error) {

	c, err := parseConfig(connStr)
	if err != nil {
		log.Fatalf("Invalid connection string: %v", err)
		return nil, err
	}

	client, err := flux.NewClient(c)
	if err != nil {
		log.Fatalf("Error while creating InfluxDB client: %v", err)
		return nil, err
	}

	p := InfluxDBPublisher{
		Config: c,
		Client: client,
	}

	return p, nil
}

// InfluxDBPublisher
type InfluxDBPublisher struct {
	Config *flux.ClientConfig
	Client *flux.Client
}

// String
func (m *InfluxDBPublisher) String() string {
	return "InfluxDB Publisher"
}

// Publish
func (p InfluxDBPublisher) Publish(in <-chan *commons.Metric, err chan<- error) {
	go func() {
		for {
			select {
			case msg := <-in:

				var sendErr error
				ser := &flux.Series{
					Name:    msg.FormatFQName(),
					Columns: []string{"time", "value"},
					Points:  [][]interface{}{{msg.Timestamp.Unix(), msg.Value}},
				}

				if p.Config.IsUDP {
					sendErr = p.Client.WriteSeriesOverUDP([]*flux.Series{ser})
				} else {
					sendErr = p.Client.WriteSeries([]*flux.Series{ser})
				}

				if sendErr != nil {
					log.Printf("Error on: %v", sendErr)
					err <- sendErr
				}

			} // select
		} // for
	}() // go
}

// Finalize
func (p InfluxDBPublisher) Finalize() {
	log.Println("InfluxDB publisher is done")
}

// parseConfig parses connStr string into an InfluxDB config
//    http://user:password@127.0.0.1:8086/dbname
//    udp://user:password@127.0.0.1:4444/dbname
func parseConfig(connStr string) (*flux.ClientConfig, error) {

	u, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}

	c := &flux.ClientConfig{}
	c.IsUDP = (u.Scheme == "udp")
	c.Host = u.Host
	c.Username = u.User.Username()

	p, _ := u.User.Password()
	c.Password = p
	c.Database = strings.Replace(u.Path, "/", "", -1)

	return c, nil
}
