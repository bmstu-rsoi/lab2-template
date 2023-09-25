package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/migregal/bmstu-iu7-ds-lab2/library/core/ports/libraries"
)

type Config struct {
	HTTPAddr string `mapstructure:"http_addr"`

	Libraries libraries.Config
}

func ReadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/usr/local/etc/library")
	viper.AddConfigPath(".")

	viper.SetDefault("http_addr", ":8080")
	viper.SetDefault("libraries.user", "program")
	viper.SetDefault("libraries.password", "test")
	viper.SetDefault("libraries.database", "libraries")
	viper.SetDefault("libraries.host", "localhost")
	viper.SetDefault("libraries.port", "5432")
	viper.SetDefault("libraries.migration_interval", "10s")
	viper.SetDefault("libraries.enable_test_data", "false")

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
