package durable

import (
	"github.com/IBM/sarama"
)

var kafkaProducer sarama.SyncProducer

func KafkaConnection() sarama.SyncProducer {
	return kafkaProducer
}

func SetupKafkaProducer(broker string) error {
	if kafkaProducer != nil {
		return nil
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	var err error
	kafkaProducer, err = sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return err
	}

	return nil
}

func SetupKafkaConsumer(broker string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return sarama.NewConsumer([]string{broker}, config)
}

func CloseKafkaConsumer(consumer sarama.Consumer) error {
	if err := consumer.Close(); err != nil {
		return err
	}
	return nil
}
