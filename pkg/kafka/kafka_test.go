package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestNewKafka(t *testing.T) {
	kafkaInstance, err := NewKafka()
	if err != nil {
		t.Fatalf("Failed to create Kafka instance: %v", err)
	}

	// Проверить, что writer и reader были успешно инициализированы
	if kafkaInstance.writer == nil {
		t.Error("Writer is nil")
	}
	if kafkaInstance.reader == nil {
		t.Error("Reader is nil")
	}
}

func TestProduceMessage(t *testing.T) {
	// Создать экземпляр Kafka для тестирования
	kafkaInstance := &Kafka{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:           []string{"localhost:9092"},
			Topic:             "test-topic",
			Balancer:          nil,
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:                []string{"localhost:9092"},
			GroupID:                "",
			GroupTopics:            []string{"test-topic"},
			Topic:                  "test-topic",
		}),
	}

	topic := "test-topic"
	additionalData := &anypb.Any{
		TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
		Value:   []byte("test-value"),
	}

	err := kafkaInstance.ProduceMessage(topic, additionalData)
	if err != nil {
		t.Fatalf("Failed to produce message: %v", err)
	}

	// Проверить бы, что сообщение было успешно отправлено без ошибок
}

func TestConsumeMessage(t *testing.T) {
	// Создать экземпляр Kafka для тестирования
	kafkaInstance := &Kafka{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:           []string{"localhost:9092"},
			Topic:             "test-topic",
			Balancer:          nil,
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:                []string{"localhost:9092"},
			GroupID:                "",
			GroupTopics:            []string{"test-topic"},
			Topic:                  "test-topic",
		}),
	}


	topic := "test-topic"

	// Отправить тестовое сообщение в кафку перед тестом
	err := kafkaInstance.ProduceMessage(topic, &anypb.Any{
		TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
		Value:   []byte("test-value"),
	})
	if err != nil {
		t.Fatalf("Failed to produce message: %v", err)
	}

	additionalData, err := kafkaInstance.ConsumeMessage(topic)
	if err != nil {
		t.Fatalf("Failed to consume message: %v", err)
	}

	// Проверить, что полученное сообщение соответствует ожидаемым значениям
	if string(additionalData.Value) != "test-value" {
		t.Errorf("Unexpected message value. Expected: %v, got: %v", "test-value", string(additionalData.Value))
	}
}