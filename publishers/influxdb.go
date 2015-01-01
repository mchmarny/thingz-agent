package publishers

import (
	"log"
	"net/url"
	"strings"

	flux "github.com/influxdb/influxdb/client"
	"github.com/mchmarny/thingz-commons"
	"github.com/mchmarny/thingz-commons/types"
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

// Publish
func (p InfluxDBPublisher) Publish(in <-chan *types.MetricCollection) {
	go func() {
		for {
			select {
			case msg := <-in:
				p.send(p.convert(*msg), false)
			} // select
		} // for
	}() // go
}

// Finalize
func (p InfluxDBPublisher) Finalize() {
	log.Println("InfluxDB publisher is done")
}

// convert metric collection to series
func (p *InfluxDBPublisher) convert(m types.MetricCollection) []*flux.Series {
	list := make([]*flux.Series, 0)
	for _, v := range m.Metrics {
		list = append(list, &flux.Series{
			Name:    commons.FormatMetricName(m.Source, m.Dimension, v.Metric),
			Columns: []string{"value"},
			Points:  [][]interface{}{{v.Value}},
		})
	}
	return list
}

// send series or retry if necessary
func (p *InfluxDBPublisher) send(list []*flux.Series, retry bool) {
	var sendErr error
	if p.Config.IsUDP {
		sendErr = p.Client.WriteSeriesOverUDP(list)
	} else {
		sendErr = p.Client.WriteSeries(list)
	}

	if sendErr != nil {
		log.Printf("Error on: %v - retrying: %v", sendErr, !retry)
		if !retry {
			p.send(list, true)
		}
	}
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
