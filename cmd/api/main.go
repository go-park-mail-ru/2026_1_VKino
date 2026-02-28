package main

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_VKino/cmd/api/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
