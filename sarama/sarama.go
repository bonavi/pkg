package sarama

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"

	"pkg/errors"
)

func ConvertStructToMessage(message any, topic string) (*sarama.ProducerMessage, error) {

	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	value := make(sarama.ByteEncoder, len(bytes))
	copy(value, bytes)

	return &sarama.ProducerMessage{
		Topic:     topic,
		Key:       nil,
		Value:     value,
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Now(),
	}, nil
}

type KafkaSettingsEnv struct {
	Addrs []string `env:"KAFKA_BOOTSTRAP_SERVERS"`
	Auth  bool     `env:"KAFKA_BOOTSTRAP_AUTH"`
	User  string   `env:"KAFKA_BOOTSTRAP_USER"`
	Pass  string   `env:"KAFKA_BOOTSTRAP_PASSWORD"`
}

func NewClient(conf KafkaSettingsEnv) (sarama.Client, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Metadata.AllowAutoTopicCreation = true
	config.Producer.Retry.Max = 5
	config.Producer.Return.Errors = true
	config.Net.SASL.Enable = conf.Auth
	config.Net.SASL.User = conf.User
	config.Net.SASL.Password = conf.Pass

	client, err := sarama.NewClient(conf.Addrs, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewAsyncProducer(conf KafkaSettingsEnv) (sarama.AsyncProducer, error) {

	client, err := NewClient(conf)
	if err != nil {
		return nil, err
	}

	producerKafka, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	return producerKafka, nil
}
