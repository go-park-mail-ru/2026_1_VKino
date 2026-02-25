package app

import (
	"flag"
	"github.com/go-park-mail-ru/2026_1_VKino/pkg/config"
	"log"
)

func Run() {

	configPath := flag.String("config", "", "config file path")
	flag.Parse()

	cfg := &Config{}
	err := config.LoadConfig(*configPath, &cfg)

	if err != nil {
		log.Fatalf("Unable to load config %w", err)
	} else {
		log.Printf("Server started on %s:%d", cfg.Server.Host, cfg.Server.Port) // Проверка
	}

}
