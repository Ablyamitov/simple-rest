package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var (
	Config = LoadConfig()
	once   sync.Once
)

type Configuration struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		URL string `yaml:"url"`
	} `yaml:"db"`
	Migration struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	} `yaml:"migration"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	App struct {
		Secret string `yaml:"secret"`
	} `yaml:"app"`
}

func NewConfig() *Configuration {
	return &Configuration{}
}
func LoadConfig() *Configuration {
	config := NewConfig()
	once.Do(func() {

		configPath := "./configs/config.yaml"
		if envPath, ok := os.LookupEnv("CONFIG_PATH"); ok {
			configPath = envPath
		}

		configFile, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Error reading config.yaml: %v", err)
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			log.Fatalf("Error parsing config.yaml: %v", err)
		}
	})
	return config
}
