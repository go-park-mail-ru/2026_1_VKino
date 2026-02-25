package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig(path string, cfg interface{}) error {

	//если не запускаем конкретный конфиг - используем локальный
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath("configs/")
		viper.SetConfigName("api")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("Error unmarshalling config, %s", err)
	}
	return nil
}
