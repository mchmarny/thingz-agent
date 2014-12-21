package publishers

import (
	"log"

	"github.com/mchmarny/thingz/types"
)

type ConsolePublisher struct{}

func (p ConsolePublisher) Publish(m *types.Metric) error {

	log.Printf("Publishing: %v", m)

	return nil

}
