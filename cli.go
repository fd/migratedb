package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "migratedb"
	app.Usage = "migrate database schemas"
	app.Flags = []cli.Flag{
		cli.StringFlag{"path", "migrations", "path of the migrations directory"},
	}
	app.Action = func(c *cli.Context) {
		var (
			dburl string
		)

		if len(c.Args()) != 1 {
			fmt.Printf("error: %s\n", "no database url was provided")
			os.Exit(1)
		} else {
			dburl = c.Args()[0]
		}

		migrations, err := LoadMigrations(c.String("path"))
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

	app.Run(os.Args)
}
