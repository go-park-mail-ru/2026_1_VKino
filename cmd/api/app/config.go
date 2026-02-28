package app

import "github.com/go-park-mail-ru/2026_1_VKino/pkg/server"

type Config struct {
	Server server.Config `mapstructure:"server"`
}
