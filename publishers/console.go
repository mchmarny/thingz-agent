package publishers

import (
	"fmt"

	"github.com/mchmarny/thingz-agent/types"
)

const (
	LINE = "------------------------------------------------------"
)

func NewConsolePublisher() (Publisher, error) {
	return ConsolePublisher{}, nil
}

type ConsolePublisher struct{}

func (p ConsolePublisher) Publish(m *types.MetricCollection) error {

	fmt.Println(LINE)
	fmt.Printf("Group: %s\n", m.Group)
	fmt.Println(LINE)

	for _, d := range m.Metrics {
		fmt.Printf("|%20s|%15s|%15v|\n",
			d.Timestamp.Format("2006-01-02T15:04:05"),
			d.Dimension, d.Value)
	}

	return nil

}
