package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type RabbitMQ struct {
	HostName  string `yaml:"host"`
	Port      uint16 `yaml:"port"`
	QueueName string `yaml:"queue"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type AppConfig struct {
	RabbitMQ   `yaml:"rabbit_mq"`
	HTTPConfig `yaml:"http"`
}

func Load(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return &cfg, nil
}
