package publishers

import (
	"fmt"

	"github.com/mchmarny/thingz-commons/types"
)

const (
	LINE = "------------------------------------------------------"
)

// NewConsolePublisher
func NewConsolePublisher() (Publisher, error) {
	return ConsolePublisher{}, nil
}

// ConsolePublisher
type ConsolePublisher struct{}

// Publish
func (p ConsolePublisher) Publish(in <-chan *types.MetricCollection) {

	go func() {
		for {
			select {
			case msg := <-in:

				fmt.Println(LINE)
				fmt.Printf("Source:%s Dimension:%s\n", msg.Source, msg.Dimension)
				fmt.Println(LINE)

				for _, m := range msg.Metrics {
					fmt.Printf("%20s %-15s %15v\n",
						m.Timestamp.Format("2006-01-02T15:04:05"),
						m.Metric, m.Value)
				}
			} // select
		} // for
	}() // go

}

// Finalize
func (p ConsolePublisher) Finalize() {
	fmt.Println("Console publisher is done")
}
