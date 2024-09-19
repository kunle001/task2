package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort     string `mapstructure:"server_port"`
	MonitoredDir   string `mapstructure:"monitored_directory"`
	CheckFrequency int    `mapstructure:"check_frequency"`
	APIEndpoint    string `mapstructure:"api_endpoint"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/file-mod-tracker/")
	viper.AddConfigPath("$HOME/.file-mod-tracker")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
