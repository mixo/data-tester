package main

import (
	"fmt"
	"github.com/mixo/data-tester/fixture"
	"os"
	"sort"
	"github.com/mixo/gosql"
	"github.com/urfave/cli/v2"
)

func main() {
	app := createApp()

	fmt.Println("Load fixtures")
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func createApp() *cli.App {
	app := &cli.App{
		Description: "The load_fixtures command is used to load test data into a database for further testing",
		Usage: " ",
		UsageText: "go run load_fixtures.go",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name:    "load_fixtures_db_driver",
				Aliases: []string{"dbd"},
				Usage:   "The database driver",
				EnvVars: []string{"load_fixtures_db_driver"},
			},
			&cli.StringFlag{
				Name:    "load_fixtures_db_host",
				Aliases: []string{"dbh"},
				Usage:   "The database host",
				EnvVars: []string{"load_fixtures_db_host"},
			},
			&cli.StringFlag{
				Name:    "load_fixtures_db_port",
				Aliases: []string{"dbt"},
				Usage:   "The database port",
				EnvVars: []string{"load_fixtures_db_port"},
			},
			&cli.StringFlag{
				Name:    "load_fixtures_db_user",
				Aliases: []string{"dbu"},
				Usage:   "The database user",
				EnvVars: []string{"load_fixtures_db_user"},
			},
			&cli.StringFlag{
				Name:    "load_fixtures_db_password",
				Aliases: []string{"dbp"},
				Usage:   "The database password",
				EnvVars: []string{"load_fixtures_db_password"},
			},
			&cli.StringFlag{
				Name:    "load_fixtures_db_name",
				Aliases: []string{"dbn"},
				Usage:   "The database name",
				EnvVars: []string{"load_fixtures_db_name"},
			},
		},
		Before: func(c *cli.Context) error {
			driver := c.String("load_fixtures_db_driver")
			host := c.String("load_fixtures_db_host")
			port := c.String("load_fixtures_db_port")
			user := c.String("load_fixtures_db_user")
			pass := c.String("load_fixtures_db_password")
			name := c.String("load_fixtures_db_name")
			if driver == "" || host == "" || port == "" || user == "" || pass == "" || name == "" {
				return cli.Exit("You must specify the database params", 1)
			}

			c.App.Metadata["db"] = gosql.DB{driver, host, port, user, pass, name}

			return nil
		},
		Commands: []*cli.Command{
			(fixture.DataLoader{}).GetCliCommand(),
		},
		CommandNotFound: func(c *cli.Context, command string) {
			fmt.Fprintf(c.App.Writer, "Unknown test '%s'\n", command)
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
