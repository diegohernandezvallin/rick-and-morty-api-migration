package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rick-and-morty-character-migration/producer/model"
	"github.com/rick-and-morty-character-migration/producer/publishing"
	"github.com/rick-and-morty-character-migration/producer/util"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.com/segmentio/kafka-go/snappy"
)

var kafkaPublisher KafkaPublisher

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) publishing.Publisher {
	if kafkaPublisher.writer == nil {
		dialer := &kafka.Dialer{
			Timeout:  10 * time.Second,
			ClientID: util.GetUUID(),
		}

		c := kafka.WriterConfig{
			Brokers:          brokers,
			Topic:            topic,
			Balancer:         &kafka.LeastBytes{},
			Dialer:           dialer,
			WriteTimeout:     10 * time.Second,
			ReadTimeout:      10 * time.Second,
			CompressionCodec: snappy.NewCompressionCodec(),
		}

		kafkaPublisher = KafkaPublisher{
			writer: kafka.NewWriter(c),
		}
	}

	return kafkaPublisher
}

func (kafkaPublisher KafkaPublisher) Publish(ctx context.Context, message model.Message) error {
	kafkaMessage, err := kafkaPublisher.encodeMessage(message)
	if err != nil {
		return err
	}

	log.Println("Publishing message with key: ", message.Key)
	return kafkaPublisher.writer.WriteMessages(ctx, kafkaMessage)
}

func (kafkaPublisher KafkaPublisher) encodeMessage(message model.Message) (kafka.Message, error) {
	value, err := json.Marshal(message.Payload)
	if err != nil {
		return kafka.Message{}, err
	}

	var headers []protocol.Header
	for k, v := range message.Headers {
		headers = append(headers, protocol.Header{Key: k, Value: []byte(v)})
	}

	return kafka.Message{
		Value:   value,
		Key:     []byte(message.Key),
		Headers: headers,
	}, nil
}
