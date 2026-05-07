package main

import (
	"log"
	"net/http"

	"jijin/backend/internal/api"
	"jijin/backend/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("backend listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, api.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
