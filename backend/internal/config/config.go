package config

import "os"

type Config struct {
	Addr      string
	DBPath    string
	AuthSecret string
}

func Load() *Config {
	return &Config{
		Addr:       getEnv("ADDR", ":8080"),
		DBPath:     getEnv("DB_PATH", "./data/glini.db"),
		AuthSecret: getEnv("AUTH_SECRET", "dev-secret-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
