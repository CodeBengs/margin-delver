package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	moduleLib "margin-delver/lib"

	_ "margin-delver/migrations"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

func main() {
	// Initialize config and logger
	cfg := moduleLib.NewAppConfig()
	logg := moduleLib.NewBaseLog(cfg)

	// Initialize DB
	db, err := moduleLib.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// parse flags and args
	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		flags.Usage()
		return
	}
	command := args[0]

	// Run goose command using context
	goose.SetDialect("mysql")
	if err := goose.RunContext(context.Background(), command, sqlDB, *dir, args[1:]...); err != nil {
		logg.SugarLog().Fatalf("migrate %v: %v", command, err)
	}
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
