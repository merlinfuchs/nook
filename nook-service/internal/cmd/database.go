package cmd

import (
	"fmt"

	"github.com/merlinfuchs/nook/nook-service/internal/entry/database"
	"github.com/urfave/cli/v2"
)

var databases = []string{"postgres"}

var databaseCMD cli.Command

func init() {
	migrateCommands := []*cli.Command{}
	for _, db := range databases {
		migrateCommands = append(migrateCommands, &cli.Command{
			Name:  db,
			Usage: fmt.Sprintf("Run migrations against the %s database.", db),
			Args:  true,
			Subcommands: []*cli.Command{
				{
					Name:  "up",
					Usage: "Migrate the database to the latest version.",
					Action: func(c *cli.Context) error {
						return database.RunMigrations(c.Context, db, database.DatabaseMigrationOpts{
							Operation: "up",
						})
					},
				},
				{
					Name:  "down",
					Usage: "Rollback the database to the earliest version.",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "danger",
							Usage: "Confirm that you want to run this command.",
						},
					},
					Action: func(c *cli.Context) error {
						if !c.Bool("danger") {
							return fmt.Errorf("this command is dangerous, use --danger flag to confirm")
						}

						return database.RunMigrations(c.Context, db, database.DatabaseMigrationOpts{
							Operation: "down",
						})
					},
				},
				{
					Name:  "version",
					Usage: "Print the current database version.",
					Action: func(c *cli.Context) error {
						return database.RunMigrations(c.Context, db, database.DatabaseMigrationOpts{
							Operation: "version",
						})
					},
				},
				{
					Name:  "list",
					Usage: "List all available database migrations.",
					Action: func(c *cli.Context) error {
						return database.RunMigrations(c.Context, db, database.DatabaseMigrationOpts{
							Operation: "list",
						})
					},
				},
				{
					Name:  "force",
					Usage: "Force a specific migration version.",
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:  "version",
							Usage: "The target version to force to.",
						},
						&cli.BoolFlag{
							Name:  "danger",
							Usage: "Confirm that you want to run this command.",
						},
					},
					Action: func(c *cli.Context) error {
						if !c.Bool("danger") {
							return fmt.Errorf("this command is dangerous, use --danger flag to confirm")
						}

						return database.RunMigrations(c.Context, db, database.DatabaseMigrationOpts{
							Operation:     "force",
							TargetVersion: c.Int("version"),
						})
					},
				},
				{
					Name:  "to",
					Usage: "Migrate the database to a specific version.",
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:  "version",
							Usage: "The target version to migrate to.",
						},
						&cli.BoolFlag{
							Name:  "danger",
							Usage: "Confirm that you want to run this command.",
						},
					},
					Action: func(c *cli.Context) error {
						if !c.Bool("danger") {
							return fmt.Errorf("this command is dangerous, use --danger flag to confirm")
						}

						return database.RunMigrations(c.Context, db, database.DatabaseMigrationOpts{
							Operation:     "to",
							TargetVersion: c.Int("version"),
						})
					},
				},
			},
		})
	}

	backupCommands := []*cli.Command{}
	for _, db := range databases {
		backupCommands = append(backupCommands, &cli.Command{
			Name:  db,
			Usage: fmt.Sprintf("Backup the %s database.", db),
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "Create a new backup of the database.",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Usage:    "The name of the backup.",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						return database.RunBackup(c.Context, db, database.DatabaseBackupOpts{
							Operation: "create",
							Name:      c.String("name"),
						})
					},
				},
				{
					Name:  "restore",
					Usage: "Restore a backup of the database.",
					Action: func(c *cli.Context) error {
						return database.RunBackup(c.Context, db, database.DatabaseBackupOpts{
							Operation: "restore",
						})
					},
				},
			},
		})
	}

	databaseCMD = cli.Command{
		Name:  "database",
		Usage: "Manage and migrate databases used by Nook.",
		Subcommands: []*cli.Command{
			{
				Name:        "migrate",
				Description: "Run database migrations.",
				Subcommands: migrateCommands,
			},
			{
				Name:        "backup",
				Description: "Backup the database.",
				Subcommands: backupCommands,
			},
		},
	}

}
