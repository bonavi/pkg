package sarama

import (
	"fmt"

	"github.com/Shopify/sarama"

	"pkg/errors"
	"pkg/log"
)

func ErrorHandler(producerKafka sarama.AsyncProducer, biddingStopper func() error) {

	// Слушаем канал ошибок от продюсера кафки
	for err := range producerKafka.Errors() {

		if biddingStopper != nil {

			// Вызываем функцию остановки аукциона
			if biddingStopperErr := biddingStopper(); biddingStopperErr != nil {
				log.Error(biddingStopperErr)
			}
		}

		// Кастим ошибку к типу ProducerError
		var kafkaErr *sarama.ProducerError
		if errors.As(err, &kafkaErr) && kafkaErr != nil && kafkaErr.Msg != nil {

			// Получаем дополнительные данные из ошибки
			var key, value []byte
			if kafkaErr.Msg.Key != nil {
				key, _ = kafkaErr.Msg.Key.Encode()
			}
			if kafkaErr.Msg.Value != nil {
				value, _ = kafkaErr.Msg.Value.Encode()
			}

			// Логируем ошибку
			log.Error(errors.InternalServer.Wrap(kafkaErr), log.ParamsOption(
				"topic", kafkaErr.Msg.Topic,
				"key", string(key),
				"value", string(value),
				"metadata", fmt.Sprintf("%v", kafkaErr.Msg.Metadata),
				"timestamp", kafkaErr.Msg.Timestamp.String(),
			))
		} else {
			log.Error(errors.InternalServer.Wrap(err))
		}
	}
}
