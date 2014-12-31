package publishers

import (
	"fmt"

	types "github.com/mchmarny/thingz-commons"
)

const (
	LINE = "------------------------------------------------------"
)

func NewConsolePublisher() (Publisher, error) {
	return ConsolePublisher{}, nil
}

type ConsolePublisher struct{}

func (p ConsolePublisher) Publish(in <-chan *types.MetricCollection) {

	go func() {
		for {
			select {
			case msg := <-in:

				fmt.Println(LINE)
				fmt.Printf("Group: %s\n", msg.Group)
				fmt.Println(LINE)

				for _, m := range msg.Metrics {
					fmt.Printf("%20s %-15s %15v\n",
						m.Timestamp.Format("2006-01-02T15:04:05"),
						m.Dimension, m.Value)
				}
			} // select
		} // for
	}() // go

}

// Finalize
func (p ConsolePublisher) Finalize() {
	fmt.Println("Console publisher is done")
}
