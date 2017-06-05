package tonic

import (
	"github.com/Shopify/sarama"
	"time"
)

type KafkaClass struct {
	AppName  string
	Enabled  bool
	Producer sarama.AsyncProducer
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

	producer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)
	if err != nil {
		return err
	}

	Kafka.Producer = producer

	return
}
