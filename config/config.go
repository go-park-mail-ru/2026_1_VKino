package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Заглушка для конфига БД
type DatabaseConfig struct {
}

// Заглушка для конфига JWT
type JWTConfig struct {
}

type ServerConfig struct {
	host string `yaml:"server_host"`
	port int    `yaml:"server_port"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

func (cfg *Config) LoadConfig() error {

	viper.AddConfigPath("./config/")
	viper.SetConfigName("config.dev")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Error reading config file, %s", err)
	}

	Server := ServerConfig{
		host: viper.GetString("server_host"),
		port: viper.GetInt("server_port"),
	}
	Database := DatabaseConfig{}
	JWT := JWTConfig{}

	cfg.Server = Server
	cfg.Database = Database
	cfg.JWT = JWT

	return nil
}
