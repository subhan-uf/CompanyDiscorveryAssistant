package config

import (
    "log"
    "os"
)

type Config struct {
    DatabaseURL string
    FlaskURL    string
    Port        string
}

func getenv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func Load() Config {
    cfg := Config{
        DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smartassistant?sslmode=disable"),
        FlaskURL:    getenv("FLASK_URL", "http://localhost:5000"),
        Port:        getenv("PORT", "8080"),
    }
    log.Printf("config loaded: port=%s flask=%s", cfg.Port, cfg.FlaskURL)
    return cfg
}
