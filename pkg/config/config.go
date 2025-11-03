package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Casbin   CasbinConfig
	Media    MediaConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type CasbinConfig struct {
	ModelPath string
}

type MediaConfig struct {
	ImgBBAPIKey string
	ImgBBTTL    uint64
}

type RedisConfig struct {
	URL       string
	QueueName string
}

func Load() (*Config, error) {
	// Parse media TTL (seconds)
	var imgbbTTL uint64 = 0
	if v := getEnv("IMGBB_TTL", "0"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			imgbbTTL = n
		}
	}
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	expirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		expirationHours = 24
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "simple_go_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours: expirationHours,
		},
		Casbin: CasbinConfig{
			ModelPath: getEnv("CASBIN_MODEL_PATH", "./configs/casbin_model.conf"),
		},
		Media: MediaConfig{
			ImgBBAPIKey: getEnv("IMGBB_API_KEY", ""),
			ImgBBTTL:    imgbbTTL,
		},
		Redis: RedisConfig{
			URL:       getEnv("REDIS_URL", "redis://localhost:6379/0"),
			QueueName: getEnv("TRANSLATION_QUEUE", "translation_jobs"),
		},
	}, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
