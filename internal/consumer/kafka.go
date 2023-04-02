package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/hablof/logistic-package-api-facade/internal/config"
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	consumer      sarama.Consumer
}

func NewKafkaConsumer(cfg config.Kafka) (*KafkaConsumer, error) {

	cg, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, sarama.NewConfig())
	if err != nil {
		return nil, err
	}

}
