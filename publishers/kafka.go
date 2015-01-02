package publishers

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	kafka "github.com/Shopify/sarama"
	"github.com/mchmarny/thingz-commons"
)

// NewKafkaPublisher factors new KafkaPublisher as Publisher
// this will create the producer up front in case we need to panic
// after that just logging errors
func NewKafkaPublisher(src, args string) (Publisher, error) {

	if len(src) < 1 || len(args) < 1 {
		log.Fatalln("Invalid arguments. Both src and args required")
		return nil, errors.New("Invalid arguments")
	}

	argParts := strings.Split(args, ",")
	if len(argParts) < 2 {
		log.Fatalln("Topic, at least one broker required: TOPIC,HOST:PORT,HOST:PORT...")
		return nil, errors.New("Invalid arguments")
	}

	topic := argParts[0]
	brokers := argParts[1:len(argParts)]
	log.Printf("Kafka - Publishing to topic:%s > %v", topic, brokers)

	// TODO: Parameterize the client configuration
	producerClient, err := kafka.NewClient(src, brokers,
		&kafka.ClientConfig{
			MetadataRetries:            3,
			WaitForElection:            3 * time.Second,
			BackgroundRefreshFrequency: 0,
		})
	if err != nil {
		log.Fatalf("Error while creating client: %v", err)
		return nil, err
	}

	// TODO: Parameterize the producer configuration
	producer, err := kafka.NewProducer(producerClient,
		&kafka.ProducerConfig{
			Partitioner:     kafka.NewHashPartitioner,
			RequiredAcks:    kafka.WaitForLocal,
			MaxMessageBytes: 1000000,
			RetryBackoff:    500 * time.Millisecond,
		})
	if err != nil {
		log.Fatalf("Error while creating producer: %v", err)
		return nil, err
	}

	return &KafkaPublisher{
		Topic:    topic,
		Producer: producer,
	}, nil

}

// KafkaPublisher represents the Kafka queue type
type KafkaPublisher struct {
	Topic    string
	Producer *kafka.Producer
}

// String
func (m *KafkaPublisher) String() string {
	return "Kafka Publisher"
}

// Publish pushes the individual series to queue
func (p KafkaPublisher) Publish(in <-chan *commons.Metric, err chan<- error) {

	go func() {
		for {
			select {
			case msg := <-in:

				data, jsonErr := msg.ToBytes()

				if jsonErr != nil {
					log.Printf("Error on marshaling metric: %v", jsonErr)
					err <- jsonErr
				}

				p.Producer.Input() <- &kafka.MessageToSend{
					Topic: p.Topic,
					Key:   nil, // this will gen a hash on server side
					Value: kafka.StringEncoder(data),
				}

			case per := <-p.Producer.Errors():
				if per != nil {
					log.Printf("Publisher error: %v", per)
					err <- errors.New(fmt.Sprintf("Kafka publisher error: %v", per))
				}
			} // select
		} // for
	}() // go

}

// Finalize
func (p KafkaPublisher) Finalize() {

	p.Producer.Close()
	log.Println("InfluxDB publisher is done")
}
