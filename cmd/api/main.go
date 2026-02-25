package main

import (
	"github.com/go-park-mail-ru/2026_1_VKino/config"
	"log"
)

func main() {
	cfg := config.Config{}
	err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to load config %w", err)
	} else {
		log.Println("Config loaded successfully")
	}

}
