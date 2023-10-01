package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	server "pvpc-backend/internal/platform/http"
	"pvpc-backend/pkg/logger"
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
	// REE API configuration
	RedataApiUrl  string `split_words:"true" required:"true"`
	EsiosApiUrl   string `split_words:"true" required:"true"`
	EsiosApiToken string `split_words:"true" required:"true"`
}

func main() {
	var err error

	cfg := loadConfig()
	configureLogger(cfg.LogLevel)

	db, err := databaseConnection(cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbTimeout)
	if err != nil {
		logger.Fatal("Error connecting to database", "err", err)
	}
	logger.Debug("Database connection established")
	defer db.Close()

	srv := server.NewHttpServer(cfg.Host, cfg.Port, cfg.Env, cfg.ShutdownTimeout, db, cfg.DbTimeout, cfg.RedataApiUrl, cfg.EsiosApiUrl, cfg.EsiosApiToken)
	srv.Run()
}

func configureLogger(level string) {
	loggerOpts := &slog.HandlerOptions{Level: logger.ParseLevel(level)}
	logger.SetDefaultLoggerJSON(loggerOpts)
}

func loadConfig() config {
	logger.Debug("Loading config...")
	var cfg config
	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file", "err", err)
	}
	if err := envconfig.Process("PVPC", &cfg); err != nil {
		logger.Fatal("Error processing env config", "err", err)
	}
	logger.Debug("Config loaded", "config", cfg)
	return cfg
}

func databaseConnection(user, pass, host string, port uint, name string, timeout time.Duration) (*sql.DB, error) {
	logger.Debug("Connecting to database...")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?connect_timeout=%d", user, pass, host, port, name, timeout)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	logger.Debug("Testing database connection...")
	err = db.Ping()
	return db, err
}
