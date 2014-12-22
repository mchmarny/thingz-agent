package publishers

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	flux "github.com/influxdb/influxdb/client"
	"github.com/mchmarny/thingz/types"
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

type InfluxDBPublisher struct {
	Config *flux.ClientConfig
	Client *flux.Client
}

func (p InfluxDBPublisher) Publish(m *types.MetricCollection) error {

	keys := make([]string, 0, len(m.Metrics))
	vals := make([]interface{}, 0, len(m.Metrics))

	for _, v := range m.Metrics {
		keys = append(keys, v.Dimension)
		vals = append(vals, v.Value)
	}

	s := &flux.Series{
		Name:    m.Group,
		Columns: keys,
		Points:  [][]interface{}{vals},
	}

	var sendErr error
	sendData := []*flux.Series{s}
	if p.Config.IsUDP {
		sendErr = p.Client.WriteSeriesOverUDP(sendData)
	} else {
		sendErr = p.Client.WriteSeries(sendData)
	}

	if sendErr != nil {
		log.Printf("H:%s U:%s P:%s D:%s",
			p.Config.Host, p.Config.Username, p.Config.Password, p.Config.Database)
		log.Fatalf("Error on series write: %v", sendErr)
		return sendErr
	} else {
		return nil
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

func printQueries() {

	fmt.Println("select total, user, sys from cpu group by time(10s) where time > now() - 30m limit 1000;")

}
