package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/dyliu/tonic/configs"
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

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	cfg.Producer.Compression = sarama.CompressionNone     // Compress messages
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, cfg)
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
