package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
)

const usage = `migratedb - migrate database schemas

Usage:
  migratedb [--path=<dir>] <database_url>
  migratedb -h | --help
  migratedb --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --path=<dir>  Path to migrations directory [default: migrations].
`

func main() {
	var (
		args, _ = docopt.Parse(usage, nil, true, "1.0", false)
		dburl   = args["<database_url>"].(string)
		path    = args["--path"].(string)
	)

	migrations, err := LoadMigrations(path)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	adapter, err := Open(dburl)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	migrator := &Migrator{migrations, adapter, log.New(os.Stdout, "", 0)}
	migrator.Run()
}
