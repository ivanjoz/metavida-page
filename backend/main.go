package main

import (
	"log"

	"metavida/backend/config"
	"metavida/backend/httpapi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server, err := httpapi.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Init(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Metavida backend listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
