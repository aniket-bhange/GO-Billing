package config

import (
	"os"
)

type AMQConfig struct {
	Url string
}

type FireBaseConfig struct {
	FilePath    string
	DatabaseUrl string
}

type DatabaseConfig struct {
	DB_USER   string
	DB_PWD    string
	DB_HOST   string
	DB_PORT   int
	DB_NAME   string
	DB_DRIVER string
}

type Config struct {
	AMQ      AMQConfig
	Firebase FireBaseConfig
	DB       DatabaseConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		AMQ: AMQConfig{
			Url: getEnv("AMQ_URL", ""),
		},
		Firebase: FireBaseConfig{
			FilePath:    getEnv("KEY_PATH", ""),
			DatabaseUrl: getEnv("DB_URL", ""),
		},
		DB: DatabaseConfig{
			DB_USER:   "root",
			DB_PORT:   3306,
			DB_HOST:   "127.0.0.1",
			DB_PWD:    "admin",
			DB_NAME:   "sample",
			DB_DRIVER: "mysql",
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
