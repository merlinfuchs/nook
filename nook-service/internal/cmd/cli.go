package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/merlinfuchs/nook/nook-service/internal/entry/service"
	"github.com/urfave/cli/v2"
)

var CLI = cli.App{
	Name:        "nook",
	Description: "Nook CLI",
	Commands: []*cli.Command{
		{
			Name:  "service",
			Usage: "Start the Nook service.",
			Action: func(c *cli.Context) error {
				ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
				defer cancel()

				return service.RunService(ctx)
			},
		},
		&databaseCMD,
	},
}

func Execute() {
	if err := CLI.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
