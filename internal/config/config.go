package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	WeatherAPIKey string
	WeatherCity   string
	GitHubUser    string
	GitHubToken   string
	CryptoCoins   []string
}

func Load() Config {
	godotenv.Load()

	coins := strings.Split(os.Getenv("CRYPTO_COINS"), ",")
	if len(coins) == 1 && coins[0] == "" {
		coins = []string{"bitcoin", "ethereum", "solana"}
	}

	return Config{
		WeatherAPIKey: os.Getenv("OPENWEATHER_API_KEY"),
		WeatherCity:   envDefault("WEATHER_CITY", "Istanbul"),
		GitHubUser:    envDefault("GITHUB_USERNAME", "TRINITY-21"),
		GitHubToken:   os.Getenv("GITHUB_TOKEN"),
		CryptoCoins:   coins,
	}
}

func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
