package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	defaultHTTPPort       = "8080"
	defaultAccessTokenTTL = 7 * 24 * time.Hour // 1 week
	defaultMigrationsPath = "file://migrations"
)

type (
	Config struct {
		HTTP   HTTPConfig
		Logger LoggerConfig
		DB     DatabaseConfig
		Auth   AuthConfig
	}

	HTTPConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	LoggerConfig struct {
		LoggerEnv string
	}

	DatabaseConfig struct {
		Host           string
		Port           string
		User           string
		Password       string
		DBName         string
		SSLMode        string
		DSN            string
		MigrationsPath string `mapstructure:"migrationsPath"`
	}

	JWTConfig struct {
		AccessTokenTTL time.Duration `mapstructure:"accessTokenTTL"`
		SigningKey     string
	}

	AuthConfig struct {
		JWT JWTConfig
	}
)

func Init(configDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshalConfig(&cfg); err != nil {
		return nil, err
	}

	setFormEnv(&cfg)

	cfg.DB.DSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	return &cfg, nil
}

func unmarshalConfig(cfg *Config) error {
	if err := viper.UnmarshalKey("http_server", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("db", &cfg.DB); err != nil {
		return err
	}
	return nil
}

func parseConfigFile(configDir string) error {
	viper.SetConfigName("main")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func setFormEnv(cfg *Config) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	cfg.Logger.LoggerEnv = os.Getenv("LOGG_ENV")

	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.DBName = os.Getenv("DB_NAME")
	cfg.DB.SSLMode = os.Getenv("DB_SSLMODE")

	cfg.Auth.JWT.SigningKey = os.Getenv("SIGNING_KEY")
}

func populateDefaults() {
	viper.SetDefault("http_server.port", defaultHTTPPort)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("db.migrationsPath", defaultMigrationsPath)
}
