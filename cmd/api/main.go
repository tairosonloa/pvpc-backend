package main

import (
	"log"
	"time"

	"go-pvpc/internal/platform/server"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	// Server configuration
	Host            string        `split_words:"true" default:"localhost"`
	Port            uint          `split_words:"true" default:"8080"`
	ShutdownTimeout time.Duration `split_words:"true" default:"10s"`
	Env             string        `split_words:"true" required:"true"`
	// Database configuration
	DbUser    string        `split_words:"true" default:"test_db_user"`
	DbPass    string        `split_words:"true" default:"test_db_pass"`
	DbHost    string        `split_words:"true" default:"localhost"`
	DbPort    uint          `split_words:"true" default:"5432"`
	DbName    string        `split_words:"true" default:"test_db_name"`
	DbTimeout time.Duration `split_words:"true" default:"5s"`
}

func main() {
	cfg := load_config()
	srv := server.New(cfg.Host, cfg.Port, cfg.Env, cfg.ShutdownTimeout)
	srv.Run()
}

func load_config() config {
	var cfg config
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := envconfig.Process("PVPC", &cfg); err != nil {
		log.Fatal("Error processing env config", err)
	}
	return cfg
}
