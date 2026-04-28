package sarama

import (
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
)

type CreateTopicRequest struct {
	Topic             string
	Partitions        int32
	ReplicationFactor int16
}

// CreateTopics creates topics in the kafka.
func CreateTopics(conf KafkaSettingsEnv, reqs ...CreateTopicRequest) error {

	client, err := NewClient(conf)
	if err != nil {
		return fmt.Errorf("migrate: init kafka client %w", err)
	}

	adm, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return fmt.Errorf("migrate: create cluster admin: %w", err)
	}

	defer func() { _ = adm.Close() }()

	for _, req := range reqs {
		details := sarama.TopicDetail{
			NumPartitions:     req.Partitions,
			ReplicationFactor: req.ReplicationFactor,
			ReplicaAssignment: nil,
			ConfigEntries:     nil,
		}
		
		err = adm.CreateTopic(req.Topic, &details, false)
		if err != nil {
			if errors.Is(err, sarama.ErrTopicAlreadyExists) {
				// Пробуем увеличить партиции если нужно
				err = adm.CreatePartitions(req.Topic, req.Partitions, nil, false)
				if err != nil {
					// Если партиций уже столько же или больше — не ошибка
					if errors.Is(err, sarama.ErrInvalidPartitions) {
						continue
					}
					return fmt.Errorf("migrate: alter partitions %q: %w", req.Topic, err)
				}
				continue
			}
			return fmt.Errorf("migrate: create topic %q: %w", req.Topic, err)
		}
	}

	return nil
}
