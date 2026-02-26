package app

import (
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2026_1_VKino/pkg/config"
	"log"
)

func Run() error {

	configPath := flag.String("config", "", "config file path")
	flag.Parse()

	cfg := &Config{}
	err := config.LoadConfig(*configPath, &cfg)

	if err != nil {
		return fmt.Errorf("Unable to load config %w", err)
	} else {
		log.Printf("Server started on %s:%d", cfg.Server.Port)
	}
	return nil
}
