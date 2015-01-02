package publishers

import (
	"fmt"

	"github.com/mchmarny/thingz-commons"
)

// NewConsolePublisher
func NewConsolePublisher() (Publisher, error) {
	return ConsolePublisher{}, nil
}

// ConsolePublisher
type ConsolePublisher struct{}

// Publish
func (p ConsolePublisher) Publish(in <-chan *commons.Metric, err chan<- error) {

	go func() {
		for {
			select {
			case msg := <-in:
				fmt.Println(msg)
			} // select
		} // for
	}() // go

}

// Finalize
func (p ConsolePublisher) Finalize() {
	fmt.Println("Console publisher is done")
}

// String
func (m *ConsolePublisher) String() string {
	return "Console Publisher"
}
