package config

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration represents the main application configuration
type Configuration struct {
	Server ServerConfig `json:"server"`
	DB     DBConfig     `json:"db"`
	Redis  RedisConfig  `json:"redis"`
	App    AppConfig    `json:"app"`
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port int `json:"port"`
}

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	MessageBatchSize    int    `json:"messageBatchSize"`
	WebhookURL          string `json:"webhookUrl"`
	MaxContentLength    int    `json:"maxContentLength"`
	MessageSendDryRun   bool   `json:"messageSendDryRun"`
	MessageSendInterval int    `json:"messageSendInterval"`
}

// LoadConfig loads the configuration from config.json or returns the default configuration
func LoadConfig() *Configuration {
	// Default configuration
	config := &Configuration{
		Server: ServerConfig{
			Port: 8080,
		},
		DB: DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "messaging",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		App: AppConfig{
			MessageBatchSize:    2,
			WebhookURL:          "https://webhook.site/",
			MaxContentLength:    1000,
			MessageSendDryRun:   false,
			MessageSendInterval: 2,
		},
	}

	// Try to load configuration from the file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Printf("Warning: Could not open config file: %v. Using default configuration.", err)
		return config
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(config)
	if err != nil {
		log.Printf("Warning: Could not parse config file: %v. Using default configuration.", err)
		return config
	}

	return config
}
