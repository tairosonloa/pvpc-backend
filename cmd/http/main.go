package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	server "pvpc-backend/internal/platform/http"
)

type config struct {
	// Server configuration
	Host            string        `split_words:"true" default:"localhost"`
	Port            uint          `split_words:"true" default:"8080"`
	ShutdownTimeout time.Duration `split_words:"true" default:"10s"`
	Env             string        `split_words:"true" required:"true"`
	LogLevel        string        `split_words:"true" default:"info"`
	// Database configuration
	DbUser    string        `split_words:"true" default:"test_db_user"`
	DbPass    string        `split_words:"true" default:"test_db_pass"`
	DbHost    string        `split_words:"true" default:"localhost"`
	DbPort    uint          `split_words:"true" default:"5432"`
	DbName    string        `split_words:"true" default:"test_db_name"`
	DbTimeout time.Duration `split_words:"true" default:"5s"`
}

func main() {
	var err error

	cfg := loadConfig()
	configureLogger(cfg.LogLevel)

	db, err := databaseConnection(cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbTimeout)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	srv := server.NewHttpServer(cfg.Host, cfg.Port, cfg.Env, cfg.ShutdownTimeout, db, cfg.DbTimeout)
	srv.Run()
}

func configureLogger(level string) {
	log.SetDefault(log.NewWithOptions(os.Stderr, log.Options{
		Level:           log.ParseLevel(level),
		Prefix:          "pvpc",
		ReportTimestamp: true,
		TimeFunction:    log.NowUTC,
		TimeFormat:      "2006/01/02T15:04:05Z",
	}))
}

func loadConfig() config {
	var cfg config
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: %v", err)
	}
	if err := envconfig.Process("PVPC", &cfg); err != nil {
		log.Fatal("Error processing env config: %v", err)
	}
	return cfg
}

func databaseConnection(user, pass, host string, port uint, name string, timeout time.Duration) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?connect_timeout=%d", user, pass, host, port, name, timeout)
	return sql.Open("pgx", connStr)
}
