package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/caarlos0/env"

	_ "github.com/joho/godotenv/autoload" // load env vars
)

var (
	cfg config
)

type config struct {
	Compress bool   `env:"COMPRESS"`
	Debug    bool   `env:"DEBUG"`
	Port     string `env:"PORT"`
}

func initialize() {
	// NOTE: For development, github.com/joho/godotenv/autoload
	// loads env variables from .env file for you.

	// Read configuration from env variables
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// log environment variables
	if cfg.Debug {
		prettyCfg, _ := json.MarshalIndent(cfg, "", "  ")
		log.Printf("Configuration: \n%v", string(prettyCfg))
	}
}
