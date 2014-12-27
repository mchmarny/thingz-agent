package publishers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	flux "github.com/influxdb/influxdb/client"
	"github.com/mchmarny/thingz-agent/types"
)

// NewInfluxDBPublisher parses connection string to InfluxDB
// and returned a configured version of the publisher
func NewInfluxDBPublisher(src, connStr string, verbose bool) (Publisher, error) {
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
		Source:  src,
		Config:  c,
		Client:  client,
		Verbose: verbose,
	}

	return p, nil
}

type InfluxDBPublisher struct {
	Source  string
	Config  *flux.ClientConfig
	Client  *flux.Client
	Verbose bool
}

func (p InfluxDBPublisher) Publish(m *types.MetricCollection) {
	list := make([]*flux.Series, 0)
	for _, v := range m.Metrics {
		list = append(list, &flux.Series{
			Name:    fmt.Sprintf("src.%s.dim.%s.met.%s", p.Source, m.Group, v.Dimension),
			Columns: []string{"value"},
			Points:  [][]interface{}{{v.Value}},
		})
	}
	go p.send(list, false)
}

func (p InfluxDBPublisher) send(list []*flux.Series, retry bool) {
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

	c.HttpClient = http.DefaultClient
	c.HttpClient.Timeout = 5 * time.Second

	return c, nil
}
