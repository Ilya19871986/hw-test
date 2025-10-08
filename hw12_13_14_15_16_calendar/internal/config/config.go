package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logger  LoggerConfig  `yaml:"logger"`
	Storage StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type StorageConfig struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
