package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    int
	Environment             string
	MaxWorkers              int
	GracefulShutDownTimeout int
	VaultToken              string
	SessionTimeout          int

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDB       string

	FrontBasePath  string
	SearchAffinity int
}

func LoadConfig() (*Config, error) {
	godotenv.Load()
	port, _ := strconv.Atoi(getEnvOrDefault("SMARAKA_PORT", "1323"))
	maxWorkers, _ := strconv.Atoi(getEnvOrDefault("SMARAKA_MAX_WORKERS", "10"))
	shutdownTimeout, _ := strconv.Atoi(getEnvOrDefault("SMARAKA_GRACE_TIMEOUT", "60"))
	dbPort, _ := strconv.Atoi(requireEnv("SMARAKA_PG_PORT"))
	sessionTimeout, _ := strconv.Atoi(getEnvOrDefault("SMARAKA_TIMEOUT_MINUTES", "262800"))
	searchAffinity, _ := strconv.Atoi(getEnvOrDefault("SMARAKA_SEARCH_AFFINITY", "60"))

	config := &Config{
		Port:                    port,
		Environment:             getEnvOrDefault("SMARAKA_ENV", "development"),
		MaxWorkers:              maxWorkers,
		GracefulShutDownTimeout: shutdownTimeout,
		VaultToken:              requireEnv("SMARAKA_VAULT_TOKEN"),
		SessionTimeout:          sessionTimeout,

		PostgresUser:     requireEnv("SMARAKA_PG_USER"),
		PostgresPassword: requireEnv("SMARAKA_PG_PASS"),
		PostgresHost:     requireEnv("SMARAKA_PG_HOST"),
		PostgresPort:     dbPort,
		PostgresDB:       requireEnv("SMARAKA_PG_DB"),

		FrontBasePath:  getEnvOrDefault("PUBLIC_SMARAKA_FRONT_BASE", "/app"),
		SearchAffinity: searchAffinity,
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}
func (c *Config) GetPostgresURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}
