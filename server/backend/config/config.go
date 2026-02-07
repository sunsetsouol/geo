package config

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	LLM struct {
		DashScopeAPIKey string `yaml:"dashscope_api_key"`
	} `yaml:"llm"`
}

var AppConfig Config

func InitConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}
