package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type ServiceConfig struct {
	UnpaidOrderTtl                     time.Duration `yaml:"unpaidOrderTtl"`
	UnpaidOrdersCancellingWorkersCount int           `yaml:"unpaidOrdersCancellingWorkersCount"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Service  ServiceConfig  `yaml:"service"`
	Postgres PostgresConfig `yaml:"postgres"`
}

var Data = Config{}

func Load(filename string) error {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "Failed to open config file")
	}

	err = yaml.Unmarshal(rawData, &Data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse config file")
	}
	return nil
}
