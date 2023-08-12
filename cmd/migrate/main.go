package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"
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

	cfg := load_config()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
	db, err := goose.OpenDBWithDriver("pgx", connStr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	if err := goose.Run(command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}
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
