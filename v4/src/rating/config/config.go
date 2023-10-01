package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/migregal/bmstu-iu7-ds-lab2/rating/core/ports/ratings"
)

type Config struct {
	HTTPAddr string `mapstructure:"http_addr"`

	Ratings ratings.Config
}

func ReadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/usr/local/etc/rating")
	viper.AddConfigPath(".")

	viper.SetDefault("http_addr", ":8050")
	viper.SetDefault("ratings.user", "program")
	viper.SetDefault("ratings.password", "test")
	viper.SetDefault("ratings.database", "ratings")
	viper.SetDefault("ratings.host", "localhost")
	viper.SetDefault("ratings.port", "5432")
	viper.SetDefault("ratings.migration_interval", "10s")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
