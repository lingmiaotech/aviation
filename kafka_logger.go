package tonic

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaHook struct {
	// topic
	topic string

	// Log levels allowed
	levels []logrus.Level

	// Log entry formatter
	formatter logrus.Formatter
}

func NewKafkaHook(topic string, levels []logrus.Level, formatter logrus.Formatter) (*KafkaHook, error) {

	if Kafka.Producer == nil {
		return nil, errors.New("tonic_error.kafka_not_enabled")
	}

	hook := &KafkaHook{
		topic,
		levels,
		formatter,
	}

	return hook, nil
}

func (hook *KafkaHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *KafkaHook) Fire(entry *logrus.Entry) error {

	go func(entry *logrus.Entry) {

		// Check time for partition key
		var partitionKey sarama.ByteEncoder

		// Get field time
		t, _ := entry.Data["time"].(time.Time)

		// Convert it to bytes
		b, err := t.MarshalBinary()

		if err != nil {
			fmt.Println(err)
		}

		partitionKey = sarama.ByteEncoder(b)
		// Format before writing
		b, err = hook.formatter.Format(entry)

		if err != nil {
			fmt.Println(err)
		}

		value := sarama.ByteEncoder(b)

		Kafka.Producer.Input() <- &sarama.ProducerMessage{
			Key:   partitionKey,
			Topic: hook.topic,
			Value: value,
		}

	}(entry)

	return nil
}
