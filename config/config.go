package configs

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Postgres PostgresConfig `yaml:"postgres"`
		HTTP     HTTPConfig     `yaml:"http"`
		Log      LogConfig      `yaml:"log"`
		Jwt      JwtConfig      `yaml:"jwt"`
		Kafka    KafkaConfig    `yaml:"kafka"`
	}

	KafkaConfig struct {
		Brokers   []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
		Topic     string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"orders.v1.events"`
		OrderPath string   `yaml:"path" env:"ORDERS_SERVICE_URL"`
	}
	InternalConfig struct {
		Path               string `yaml:"path" env:"API_PATH"`
		Environment        string `yaml:"environment" env:"ENVIRONMENT"`
		LogInputParamOnErr bool   `yaml:"logInputParamOnErr" env:"LOG_INPUT_PARAM_ON_ERR"`
	}

	JwtConfig struct {
		JwtSecret  string        `env:"JWT_SECRET"`
		AccessTTL  time.Duration `yaml:"accessTtl"`
		RefreshTTL time.Duration `yaml:"refreshTtl"`
	}

	CacheConfig struct {
		Url      string `yaml:"url" env:"CACHE_URL"`
		Username string `yaml:"username" env:"CACHE_USERNAME"`
		Password string `yaml:"password" env:"CACHE_PASSWORD"`
	}

	ResponseConfig struct {
		ExportError bool `yaml:"exportError" env:"RESPONSE_EXPORT_ERROR"`
	}

	LogConfig struct {
		Level              string `yaml:"level" env:"LOG_LEVEL"`
		RequestLogEnabled  bool   `yaml:"requestLogEnabled" env:"LOG_REQUEST_ENABLED"`
		RequestLogWithBody bool   `yaml:"requestLogWithBody" env:"LOG_REQUEST_WITH_BODY"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		Username string `yaml:"username" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Schema   string `yaml:"schema" env:"DB_SCHEMA"`
		DBName   string `yaml:"dbname" env:"DB_NAME"`
		SSLMode  string `yaml:"sslmode" env:"DB_SSL_MODE" env-default:"disable"`
	}

	HTTPConfig struct {
		Port         string        `yaml:"port" env:"HTTP_PORT"`
		InternalPort string        `yaml:"internal_port" env:"INTERNAL_PORT"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	}
)

func NewConfig(configDir string, envPath string) (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load(envPath); err != nil {

		fmt.Printf("Warning: %s file not found\n", envPath)
	}

	if err := cleanenv.ReadConfig(configDir, cfg); err != nil {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return cfg, nil
}

func (c *Config) GetDBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		c.Postgres.Username, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DBName, c.Postgres.SSLMode, c.Postgres.Schema)
}
