package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hablof/logistic-package-api-facade/internal/config"
	"github.com/hablof/logistic-package-api-facade/internal/kafka"
	mw "github.com/hablof/logistic-package-api-facade/internal/messagewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Error().Err(err)
		return
	}
	log.Info().Msg("config read")
	cfg := config.GetConfigInstance()

	pbChannel := make(chan []byte, cfg.Capacity)
	jsonChannel := make(chan []byte, cfg.Capacity)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kc, err := kafka.NewKafkaConsumer(cfg, pbChannel, jsonChannel)
	if err != nil {
		log.Error().Err(err).Msg("failed to create kafka consumer")
		return
	}
	log.Info().Msg("kafka consumer created")

	if err := kc.Start(ctx); err != nil {
		log.Error().Err(err).Msg("failed to start kafka consumer")
		return
	}
	log.Info().Msg("kafka consumer started")

	pbw := mw.NewPbWriter(pbChannel)
	pbw.Start(ctx)
	jw := mw.NewJsonWriter(jsonChannel)
	jw.Start(ctx)
	log.Info().Msg("writer started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	log.Info().Msg("main gorutine waiting for interrupt signal")

	<-interrupt
	cancel()
	log.Info().Msg("main gorutine recived interrupt signal")
}
