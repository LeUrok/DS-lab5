package config

import (
	"os"
)

type Config struct {
	Port       string
	CarsURL    string
	RentalURL  string
	PaymentURL string

	AuthDomain       string
	AuthClientID     string
	AuthClientSecret string
	AuthAudience     string
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8050"),
		CarsURL:    getEnv("CARS_URL", "http://localhost:8070"),
		RentalURL:  getEnv("RENTAL_URL", "http://localhost:8060"),
		PaymentURL: getEnv("PAYMENT_URL", "http://localhost:8050"),

		AuthDomain:       getEnv("AUTH_DOMAIN", ""),
		AuthClientID:     getEnv("AUTH_CLIENT_ID", ""),
		AuthClientSecret: getEnv("AUTH_CLIENT_SECRET", ""),
		AuthAudience:     getEnv("AUTH_AUDIENCE", ""),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
