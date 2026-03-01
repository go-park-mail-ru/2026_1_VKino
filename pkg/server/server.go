package server

import (
	"time"
	"net/http"
)

type Config struct {
	Port int `mapstructure:"port"`
}

func RunServer(addr string, handler http.Handler) error {
	server := http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5  * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return server.ListenAndServe()
}