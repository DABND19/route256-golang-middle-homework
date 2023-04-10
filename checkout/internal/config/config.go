package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Address        string `yaml:"address"`
	MetricsAddress string `yaml:"metricsAddress"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type LomsServiceConfig struct {
	Url string `yaml:"url"`
}

type ProductServiceConfig struct {
	Url                   string `yaml:"url"`
	AccessToken           string `yaml:"accessToken"`
	RateLimit             uint16 `yaml:"rateLimit"`
	MaxConcurrentRequests int    `yaml:"maxConcurrentRequests"`
}

type Config struct {
	Server           ServerConfig   `yaml:"server"`
	Postgres         PostgresConfig `yaml:"postgres"`
	ExternalServices struct {
		Loms    LomsServiceConfig    `yaml:"loms"`
		Product ProductServiceConfig `yaml:"product"`
	} `yaml:"externalServices"`
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
