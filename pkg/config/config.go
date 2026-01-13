package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment   string `yaml:"environment"`
	HTTPPort      string `yaml:"http_port"`
	GRPCPort      string `yaml:"grpc_port"`
	AuthSecret    string `yaml:"auth_secret"`
	DatabaseURI   string `yaml:"database_uri"`
	RedisURI      string `yaml:"redis_uri"`
	RedisPassword string `yaml:"redis_password"`
	RedisDB       int    `yaml:"redis_db"`
}

var (
	config *Config
	once   sync.Once
)

func Load(configPath string) *Config {
	once.Do(func() {
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}

		config = &Config{}
		if err := yaml.Unmarshal(data, config); err != nil {
			log.Fatalf("Failed to parse config file: %v", err)
		}

		log.Println("Configuration loaded successfully")
	})

	return config
}

func GetConfig() *Config {
	if config == nil {
		log.Fatal("Config not loaded. Call load() first")
	}

	return config
}
