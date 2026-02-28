package app

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_VKino/pkg/config"
)

func Run() error {
	configPath := flag.String("config", "", "config file path")

	flag.Parse()

	cfg := &Config{}

	err := config.LoadConfig(*configPath, &cfg)
	if err != nil {
		return fmt.Errorf("unable to load config %w", err)
	}

	log.Printf("Server started on %d", cfg.Server.Port)

	return nil
}
