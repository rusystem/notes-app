package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DB  Postgres
	RDB Redis

	Key Keys

	Cache struct {
		Ttl time.Duration `mapstructure:"ttl"`
	} `mapstructure:"cache"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Auth struct {
		SessionTTL time.Duration `mapstructure:"session_ttl"`
	} `mapstructure:"auth"`
}

type Postgres struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

type Redis struct {
	Host     string
	Port     int
	Password string
}

type Keys struct {
	Salt string
}

func New(folder, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	if err := envconfig.Process("rdb", &cfg.RDB); err != nil {
		return nil, err
	}

	if err := envconfig.Process("key", &cfg.Key); err != nil {
		return nil, err
	}

	return cfg, nil
}
