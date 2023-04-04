package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/hablof/logistic-package-api-facade/internal/config"
	"github.com/rs/zerolog/log"
)

type consumerHandler struct {
	output chan []byte
}

// Cleanup implements sarama.ConsumerGroupHandler
func (*consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler
func (ch *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Debug().Msgf("consumerHandler.ConsumeClaim: message with offset %d recived", message.Offset)
		session.MarkMessage(message, "")

		ch.output <- message.Value
	}

	return nil
}

// Setup implements sarama.ConsumerGroupHandler
func (*consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

var _ sarama.ConsumerGroupHandler = &consumerHandler{}

type KafkaConsumer struct {
	consumerGroup   sarama.ConsumerGroup
	topic           string
	consumerHandler *consumerHandler
}

func (kc *KafkaConsumer) GetChannel() <-chan []byte {
	return kc.consumerHandler.output
}

func NewKafkaConsumer(cfg config.Kafka) (*KafkaConsumer, error) {

	saramaConf := sarama.NewConfig()
	saramaConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConf.Consumer.Offsets.Initial = sarama.OffsetOldest
	log.Debug().Msg("NewKafkaConsumer(): sarama config ready")

	cg, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConf)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("NewKafkaConsumer(): consumerGroup created")

	consumerHandler := &consumerHandler{
		output: make(chan []byte, cfg.Capacity),
	}
	log.Debug().Msg("NewKafkaConsumer(): consumerGroup created")

	kc := &KafkaConsumer{
		consumerGroup:   cg,
		topic:           cfg.Topic,
		consumerHandler: consumerHandler,
	}
	log.Debug().Msg("NewKafkaConsumer(): consumer created")

	return kc, nil
}

func (kc *KafkaConsumer) Start() error {
	ctx := context.Background()
	err := kc.subscribe(ctx)

	return err
}

func (kc *KafkaConsumer) subscribe(ctx context.Context) error {
	go func() {
		for {
			if err := kc.consumerGroup.Consume(ctx, []string{kc.topic}, kc.consumerHandler); err != nil {
				log.Error().Err(err).Msg("kc.consumerGroup.Consume")
			}
		}
	}()
	return nil
}
