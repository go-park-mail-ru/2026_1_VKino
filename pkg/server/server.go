package server

import (
	"fmt"
	"net/http"

	authapp "github.com/go-park-mail-ru/2026_1_VKino/internal/app/auth"
)

type Config struct {
	Port int `mapstructure:"port"`
}

func RunServer(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Addr:", addr, "URL:", r.URL.String())
	})

	// ручки на signIn + signUp
	authService := authapp.NewService()
	authHandler := authapp.NewHandler(authService)
	authHandler.RegisterRoutes(mux)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server.ListenAndServe()
}