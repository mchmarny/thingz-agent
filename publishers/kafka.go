package publishers

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mchmarny/thingz-agent/types"
)

// NewKafkaPublisher factors new KafkaPublisher as Publisher
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
	producerClient, err := sarama.NewClient(src, brokers,
		&sarama.ClientConfig{
			MetadataRetries:            3,
			WaitForElection:            3 * time.Second,
			BackgroundRefreshFrequency: 0,
		})
	if err != nil {
		log.Fatalf("Error while creating client: %v", err)
		return nil, err
	}

	// TODO: Parameterize the producer configuration
	producer, err := sarama.NewProducer(producerClient,
		&sarama.ProducerConfig{
			Partitioner:     sarama.NewHashPartitioner,
			RequiredAcks:    sarama.WaitForLocal,
			MaxMessageBytes: 1000000,
			RetryBackoff:    500 * time.Millisecond,
		})
	if err != nil {
		log.Fatalf("Error while creating producer: %v", err)
		return nil, err
	}

	// TODO: Producer events won't fire, figure out better way to capture
	go func() {
		for {
			select {
			case err := <-producer.Errors():
				log.Printf("Error on message send: %v", err)
			case suc := <-producer.Successes():
				log.Printf("Message sent: %v", suc)
			}
		}
	}()

	return &KafkaPublisher{
		Topic:    topic,
		Producer: producer,
	}, nil

}

// KafkaPublisher represents the Kafka queue type
type KafkaPublisher struct {
	Topic    string
	Producer *sarama.Producer
}

// Publish pushes the individual series to queue
func (p KafkaPublisher) Publish(m *types.MetricCollection) {

	p.Producer.Input() <- &sarama.MessageToSend{
		Topic: p.Topic,
		Key:   nil,
		Value: sarama.StringEncoder(m.ToBytes()),
	}

}

// Finalize
func (p KafkaPublisher) Finalize() {

	p.Producer.Close()
	log.Println("InfluxDB publisher is done")
}
