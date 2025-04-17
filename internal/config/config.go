package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	defaultHTTPPort        = "8080"
	defaultAccessTokenTTL  = 24 * time.Hour
	defaultRefreshTokenTTL = 24 * time.Hour * 30
)

type (
	Config struct {
		HTTP HTTPConfig
		DB   DatabaseConfig
		Auth AuthConfig
	}

	HTTPConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	DatabaseConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
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

	return &cfg, nil
}

func unmarshalConfig(cfg *Config) error {
	if err := viper.UnmarshalKey("http_server", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
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
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.DBName = os.Getenv("DB_NAME")

	cfg.Auth.JWT.SigningKey = os.Getenv("SIGNING_KEY")
}

func populateDefaults() {
	viper.SetDefault("http_server.port", defaultHTTPPort)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
}
