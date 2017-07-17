package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/lingmiaotech/tonic/configs"
	"time"
)

type Class struct {
	AppName  string
	Enabled  bool
	Producer sarama.SyncProducer
}

var Instance Class

func InitKafka() (err error) {

	Instance.AppName = configs.GetString("app_name")
	Instance.Enabled = configs.GetBool("kafka.enabled")

	if !Instance.Enabled {
		return nil
	}

	brokers := configs.GetStringSlice("kafka.brokers")

	configs := sarama.NewConfig()
	configs.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	configs.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	configs.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	configs.Producer.Retry.Max = 3
	configs.Producer.Return.Errors = true
	configs.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, configs)
	if err != nil {
		return err
	}

	Instance.Producer = producer

	return
}

func SendMessage(message *sarama.ProducerMessage) (int32, int64, error) {
	if !Instance.Enabled {
		return -1, -1, nil
	}
	partition, offset, err := Instance.Producer.SendMessage(message)
	return partition, offset, err
}

func SendMessageWithRetry(message *sarama.ProducerMessage, retry int, backoff time.Duration) (int32, int64, error) {
	if !Instance.Enabled {
		return -1, -1, nil
	}
	partition, offset, err := Instance.Producer.SendMessage(message)
	if err != nil && retry > 1 {
		time.Sleep(backoff)
		return SendMessageWithRetry(message, retry-1, backoff)
	}
	return partition, offset, err
}
