package tonic

import (
	"github.com/Shopify/sarama"
	"time"
)

type KafkaClass struct {
	AppName  string
	Enabled  bool
	Producer sarama.SyncProducer
}

var Kafka KafkaClass

func InitKafka() (err error) {

	Kafka.AppName = Configs.GetString("app_name")
	Kafka.Enabled = Configs.GetBool("kafka.enabled")

	if !Kafka.Enabled {
		return nil
	}

	brokers := Configs.GetStringSlice("kafka.brokers")

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	kafkaConfig.Producer.Retry.Max = 3
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		return err
	}

	Kafka.Producer = producer

	return
}

func (k KafkaClass) SendMessage(message *sarama.ProducerMessage) (int32, int64, error) {
	if !k.Enabled {
		return -1, -1, nil
	}
	partition, offset, err := k.Producer.SendMessage(message)
	return partition, offset, err
}

func (k KafkaClass) SendMessageWithRetry(message *sarama.ProducerMessage, retry int, backoff time.Duration) (int32, int64, error) {
	if !k.Enabled {
		return -1, -1, nil
	}
	partition, offset, err := k.Producer.SendMessage(message)
	if err != nil && retry > 1 {
		time.Sleep(backoff)
		return k.SendMessageWithRetry(message, retry-1, backoff)
	}
	return partition, offset, err
}
