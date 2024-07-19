package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	IPLimit        int
	IPDuration     time.Duration
	APIKeyLimit    int
	APIKeyDuration time.Duration
}

func getenvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, will use environment variables")
	}
	ipLimit := getenvInt("IP_MAX_REQUESTS_PER_SECOND", 5)
	ipDuration := getenvInt("IP_BLOCK_DURATION_SECONDS", 300)
	apiKeyLimit := getenvInt("TOKEN_MAX_REQUESTS_PER_SECOND", 10)
	apiKeyDuration := getenvInt("TOKEN_BLOCK_DURATION_SECONDS", 300)
	return Config{
		IPLimit:        ipLimit,
		IPDuration:     time.Duration(ipDuration) * time.Second,
		APIKeyLimit:    apiKeyLimit,
		APIKeyDuration: time.Duration(apiKeyDuration) * time.Second,
	}
}
