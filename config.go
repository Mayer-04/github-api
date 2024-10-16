package main

import (
	"os"
	// Cargar automáticamente su archivo .env:
	_ "github.com/joho/godotenv/autoload"
)

type GithubConfig struct {
	ClientID     string
	ClientSecret string
}

type Config struct {
	Port        string
	RedirectURL string
	Github      GithubConfig
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:        os.Getenv("PORT"),
		RedirectURL: os.Getenv("REDIRECT_URL"),
		Github: GithubConfig{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		},
	}
	return cfg, nil
}

// func getEnv(key string) (string, error) {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		return "", fmt.Errorf("%s is not set", key)
// 	}
// 	return value, nil
// }
