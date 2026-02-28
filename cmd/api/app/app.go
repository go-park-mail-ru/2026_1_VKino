package app

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_VKino/pkg/config"
	"github.com/go-park-mail-ru/2026_1_VKino/pkg/server"
)

func Run() error {
	configPath := flag.String("config", "", "config file path")

	flag.Parse()

	cfg := &Config{}

	err := config.LoadConfig(*configPath, &cfg)
	if err != nil {
		return fmt.Errorf("unable to load config %w", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server started on %d", cfg.Server.Port)

	return server.RunServer(addr)
}
