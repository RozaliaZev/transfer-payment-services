package kafka

import (
	"context"
	"log"
	"os"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/anypb"
)

var broker = os.Getenv("KAFKA_BROKER")

type Kafka struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafka() (*Kafka, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "topic",
		Balancer: &kafka.LeastBytes{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "my-group",
		GroupTopics:   []string{"topic"},
	})

	return &Kafka{
		writer: writer,
		reader: reader,
	}, nil
}

func (k *Kafka) ProduceMessage(topic string, AdditionalData *anypb.Any) error {
	k.writer.Topic = topic

	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: AdditionalData.Value,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Check result sent to Kafka: %v", AdditionalData.Value)

	return nil
}

func (k *Kafka) ConsumeMessage(topic string) (*anypb.Any, error) {
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "my-group",
		GroupTopics: []string{topic},
	})

	m, err := k.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}

	log.Printf("Received message: %v\n", string(m.Value))

	additionalData := &anypb.Any{
		TypeUrl: "http://type.googleapis.com/google.protobuf.JSONValue",
		Value:   m.Value,
	}

	return additionalData, nil
}
