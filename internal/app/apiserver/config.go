package apiserver

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		URL string `yaml:"url"`
	} `yaml:"db"`
}

func NewConfig() *Config {
	return &Config{}
}
func LoadConfig() *Config {
	config := NewConfig()
	configFile, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error parsing config.yaml: %v", err)
	}
	return config
}
