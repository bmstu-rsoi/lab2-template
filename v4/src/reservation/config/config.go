package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/core/ports/reservations"
)

type Config struct {
	HTTPAddr string `mapstructure:"http_addr"`

	Reservations reservations.Config
}

func ReadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/usr/local/etc/reservation")
	viper.AddConfigPath(".")

	viper.SetDefault("http_addr", ":8070")
	viper.SetDefault("reservations.user", "program")
	viper.SetDefault("reservations.password", "test")
	viper.SetDefault("reservations.database", "reservations")
	viper.SetDefault("reservations.host", "localhost")
	viper.SetDefault("reservations.port", "5432")
	viper.SetDefault("reservations.migration_interval", "10s")

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
