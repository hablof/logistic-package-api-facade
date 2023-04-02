package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/hablof/logistic-package-api-facade/internal/config"
	"github.com/hablof/logistic-package-api-facade/internal/consumer"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		fmt.Println(err)
		return
	}
	kc, err := consumer.NewKafkaConsumer(config.GetConfigInstance())
	if err != nil {
		fmt.Println(err)
		return
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

}
