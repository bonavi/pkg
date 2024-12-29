package kafka

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
