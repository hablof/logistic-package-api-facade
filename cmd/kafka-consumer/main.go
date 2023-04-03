package main

import (
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

	kc, err := kafka.NewKafkaConsumer(config.GetConfigInstance())
	if err != nil {
		log.Error().Err(err).Msg("failed to create kafka consumer")
		return
	}
	log.Info().Msg("kafka consumer created")

	if err := kc.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start kafka consumer")
		return
	}
	log.Info().Msg("kafka consumer run")

	w := mw.NewWriter(kc.GetChannel())
	w.Start()
	log.Info().Msg("writer started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	log.Info().Msg("main gorutine waiting for interrupt signal")

	<-interrupt
	log.Info().Msg("main gorutine recived interrupt signal")
}
