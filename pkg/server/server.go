package server

type Config struct {
	Host string `mapstructure:"server_host"`
	Port int    `mapstructure:"server_port"`
}
