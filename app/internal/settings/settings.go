package settings

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	ContextKeyUserID    = "user_id"
	ContextKeyUserRole  = "role"
	JwtTokenExpDuration = time.Hour * 2
	AuthHeader          = "Authorization"
	ConfigPath          = "config.yaml"
)

type LoggerType string

var (
	LoggerTypeProd LoggerType = "prod"
	LoggerTypeDev  LoggerType = "dev"
)

// all general env reading and other init stuff should be here

type Config struct {
	JWT     JWTConfig     `yaml:"jwt"`
	Logging LoggingConfig `yaml:"logging"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type LoggingConfig struct {
	Type LoggerType `yaml:"type"`
}

func ReadConfig() (Config, error) {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("read %s: %w", ConfigPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal %s: %w", ConfigPath, err)
	}

	fmt.Println(cfg)

	if cfg.JWT.Secret == "" {
		return Config{}, fmt.Errorf("%s: jwt.secret is required", ConfigPath)
	}
	if cfg.Logging.Type == "" {
		return Config{}, fmt.Errorf("%s: logging.type is required", ConfigPath)
	}
	if cfg.Logging.Type != LoggerTypeDev && cfg.Logging.Type != LoggerTypeProd {
		return Config{}, fmt.Errorf("%s: invalid logging.type", ConfigPath)
	}

	return cfg, nil
}

func ReadDbConnectionConfig() string {
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		panic("no `POSTGRES_USER` env var")
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		panic("no `POSTGRES_PASSWORD` env var")
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		panic("no `POSTGRES_PORT` env var")
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		panic("no `POSTGRES_DB` env var")
	}
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		panic("no `POSTGRES_HOST` env var")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
}
