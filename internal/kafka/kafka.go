package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/hablof/logistic-package-api-facade/internal/config"
	"github.com/rs/zerolog/log"
)

type consumerHandler struct {
	databaseEventsOutput chan<- []byte
	tgbotEventsOutput    chan<- []byte
}

// Cleanup implements sarama.ConsumerGroupHandler
func (*consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler
func (ch *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Debug().Msgf(`consumerHandler.ConsumeClaim: message from topic "%s" with offset %d recived`, message.Topic, message.Offset)
		session.MarkMessage(message, "")

		// shitty hardcode
		switch message.Topic {
		case "omp-package-events":
			ch.databaseEventsOutput <- message.Value

		case "omp-tgbot-commands":
			ch.tgbotEventsOutput <- message.Value

		case "omp-tgbot-cache-events":
			ch.tgbotEventsOutput <- message.Value
		}
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
	topics          []string
	consumerHandler *consumerHandler
}

// func (kc *KafkaConsumer) GetChannel() <-chan []byte {
// 	return kc.consumerHandler.output
// }

func NewKafkaConsumer(cfg config.Kafka, packageEventsCh chan<- []byte, tgbotCommandsCh chan<- []byte) (*KafkaConsumer, error) {

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
		databaseEventsOutput: packageEventsCh,
		tgbotEventsOutput:    tgbotCommandsCh,
	}

	kc := &KafkaConsumer{
		consumerGroup:   cg,
		topics:          cfg.Topics,
		consumerHandler: consumerHandler,
	}
	log.Debug().Msg("NewKafkaConsumer(): consumer created")

	return kc, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	// ctx := context.Background()
	err := kc.subscribe(ctx)

	return err
}

func (kc *KafkaConsumer) subscribe(ctx context.Context) error {
	go func() {
		// for {
		if err := kc.consumerGroup.Consume(ctx, kc.topics, kc.consumerHandler); err != nil {
			log.Error().Err(err).Msg("kc.consumerGroup.Consume")
		}
		// }
	}()
	return nil
}
