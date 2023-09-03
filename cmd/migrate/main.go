package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"

	"pvpc-backend/pkg/logger"
)

type config struct {
	DbUser string `split_words:"true" default:"test_db_user"`
	DbPass string `split_words:"true" default:"test_db_pass"`
	DbHost string `split_words:"true" default:"localhost"`
	DbPort uint   `split_words:"true" default:"5432"`
	DbName string `split_words:"true" default:"test_db_name"`
}

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "internal/platform/storage/postgresql/migrations", "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	cfg := loadConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
	db, err := goose.OpenDBWithDriver("pgx", connStr)
	if err != nil {
		logger.Fatal("Error opening DB", "err", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal("Error closing DB", "err", err)
		}
	}()

	if err := goose.Run(command, db, *dir, args[1:]...); err != nil {
		logger.Fatal("Migrate error", "err", err, "command", command)
	}
}

func loadConfig() config {
	var cfg config
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file", "err", err)
	}
	if err := envconfig.Process("PVPC", &cfg); err != nil {
		logger.Fatal("Error processing env config", "err", err)
	}
	return cfg
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
