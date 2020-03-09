package cmd

import (
	"fmt"

	"github.com/bbckr/parcel/structs"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

func NewPullCommand(cfg *structs.Config, logger *log.Entry) *cli.Command {
	var opts structs.PullOptions
	return &cli.Command{
		Name:  "pull",
		Usage: "Pull a parcel from a git source",
		Action: func(c *cli.Context) error {
			return RunPull(c, cfg, &opts, logger)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Aliases:     []string{"f"},
				Usage:       "(optional) force pull parcel",
				Value:       false,
				Destination: &opts.Force,
			},
		},
		ArgsUsage: "[SOURCE]",
	}
}

func RunPull(c *cli.Context, cfg *structs.Config, opts *structs.PullOptions, logger *log.Entry) error {
	logger = logger.WithField("command", c.Command.Name)
	source := c.Args().First()

	engine := structs.NewEngineFromConfig(cfg)
	logger.WithField("source", source).Info("Preparing to pull parcel from source")
	if err := engine.Pull(source, opts, logger); err != nil {
		err = fmt.Errorf("Failed to pull source: %s", err)
		logger.Fatal(err)
		return err
	}
	return nil
}
