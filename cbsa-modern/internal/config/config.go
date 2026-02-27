package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	SortCode    string
	CompanyName string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://cbsa:cbsa_secret@localhost:5432/cbsa?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		SortCode:    getEnv("SORT_CODE", "987654"),
		CompanyName: getEnv("COMPANY_NAME", "CICS Bank Sample Application"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
