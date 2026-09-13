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
	RabbitMQ    `yaml:"rabbit_mq"`
	HTTPConfig  `yaml:"http"`
	Postgres    `yaml:"postgres"`
	RedisConfig `yaml:"redis"`
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

type Postgres struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	TTL      int    `yaml:"ttl"`
}
