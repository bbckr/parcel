package cmd

import (
	"fmt"
	"os"

	"github.com/bbckr/parcel/helpers"
	"github.com/bbckr/parcel/structs"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewRenderCommand(cfg *structs.Config, logger *log.Entry) *cli.Command {
	var opts structs.RenderOptions
	return &cli.Command{
		Name:  "render",
		Usage: "Render template output for a parcel",
		Action: func(c *cli.Context) error {
			return RunRender(c, cfg, &opts, logger)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "values",
				Aliases:     []string{"v"},
				Usage:       "(optional) path to values yaml file",
				Destination: &opts.ValuesFilePath,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "directory to output rendered files",
				Value:       ".",
				Destination: &opts.OutputDirectory,
			},
		},
		ArgsUsage: "[SOURCE]",
	}
}

func RunRender(c *cli.Context, cfg *structs.Config, opts *structs.RenderOptions, logger *log.Entry) error {
	logger = logger.WithField("command", c.Command.Name)
	source := c.Args().First()

	if _, err := os.Stat(opts.OutputDirectory); os.IsNotExist(err) {
		logger.Infof("Ensuring output directory: %s", opts.OutputDirectory)
		err := helpers.EnsureDirectory(opts.OutputDirectory)
		if err != nil {
			err = fmt.Errorf("Failed to ensure output directory: %s", err)
			logger.Fatal(err)
			return err
		}
	}

	engine := structs.NewEngineFromConfig(cfg)
	logger.WithField("source", source).Info("Loading parcel from source")
	parcel, err := engine.Load(source)
	if err != nil {
		err = fmt.Errorf("Failed to load parcel: %s", err)
		logger.Fatal(err)
		return err
	}

	logger.WithField("parcel", parcel.ID()).Infof("Preparing to render parcel templates")
	if err = engine.Render(parcel, opts, logger); err != nil {
		err = fmt.Errorf("Failed to render parcel: %s", err)
		logger.Fatal(err)
		return err
	}

	return nil
}
