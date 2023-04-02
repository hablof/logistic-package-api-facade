package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var cfg *Kafka

// GetConfigInstance returns service config
func GetConfigInstance() Kafka {
	if cfg != nil {
		return *cfg
	}

	return Kafka{}
}

type Kafka struct {
	Capacity uint64   `yaml:"capacity"`
	Topic    string   `yaml:"topic"`
	GroupID  string   `yaml:"groupId"`
	Brokers  []string `yaml:"brokers"`
}

func ReadConfigYML(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return err
	}

	return nil
}
