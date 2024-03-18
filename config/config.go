package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func Read(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
