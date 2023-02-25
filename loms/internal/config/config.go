package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
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
