package durable

import (
	"github.com/IBM/sarama"
)

var kafkaProducer sarama.SyncProducer

func SetupKafkaProducer(broker string) error {
	if kafkaProducer != nil {
		return nil
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	var err error
	kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return err
	}

	return nil
}

func KafkaConnection() sarama.SyncProducer {
	return kafkaProducer
}
