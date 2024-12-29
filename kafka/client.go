package kafka

import (
	"github.com/Shopify/sarama"
)

type KafkaSettingsEnv struct {
	Addresses []string `env:"KAFKA_BOOTSTRAP_SERVERS,required"`
	Auth      bool     `env:"KAFKA_BOOTSTRAP_AUTH,required"`
	User      string   `env:"KAFKA_BOOTSTRAP_USER,required"`
	Pass      string   `env:"KAFKA_BOOTSTRAP_PASSWORD,required"`
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

	client, err := sarama.NewClient(conf.Addresses, config)
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
