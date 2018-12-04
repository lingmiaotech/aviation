package logging

import (
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dyliu/tonic/kafka"
	"github.com/sirupsen/logrus"
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

	if kafka.Instance.Producer == nil {
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

	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	value := sarama.ByteEncoder(b)
	message := &sarama.ProducerMessage{
		Topic: hook.topic,
		Value: value,
	}

	_, _, err = kafka.SendMessageWithRetry(message, 3, 1*time.Second)

	return err
}
