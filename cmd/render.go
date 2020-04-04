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
		ArgsUsage: "[OWNER] [NAME] [VERSION]",
	}
}

func RunRender(c *cli.Context, cfg *structs.Config, opts *structs.RenderOptions, logger *log.Entry) error {
	logger = logger.WithField("command", c.Command.Name)
	owner := c.Args().Get(0)
	if len(owner) == 0 {
		err := fmt.Errorf("Missing first argument [OWNER]")
		logger.Fatal(err)
		return err
	}

	name := c.Args().Get(1)
	if len(name) == 0 {
		err := fmt.Errorf("Missing second argument [NAME]")
		logger.Fatal(err)
		return err
	}

	version := c.Args().Get(2)
	if len(version) == 0 {
		err := fmt.Errorf("Missing second argument [VERSION]")
		logger.Fatal(err)
		return err
	}

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
	logger.Info("Loading parcel")
	parcel, err := engine.Load(owner, name, version)
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
