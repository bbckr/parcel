package main

import (
	"fmt"
	"os"

	"github.com/bbckr/parcel/cmd"
	"github.com/bbckr/parcel/structs"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	Version = "local"
)

func main() {
	rootLogger := log.New()
	cfg := structs.NewConfig()

	if cfg.Debug {
		rootLogger.SetLevel(log.DebugLevel)
	}

	if cfg.LogLevel != "" {
		level, err := log.ParseLevel(cfg.LogLevel)
		if err != nil {
			panic(fmt.Sprintf("Invalid LOG_LEVEL: %s", cfg.LogLevel))
		}
		rootLogger.SetLevel(level)
	}

	if cfg.LogFormat == "json" {
		rootLogger.SetFormatter(&log.JSONFormatter{})
	}

	logger := rootLogger.WithField("service", "parcel")

	app := &cli.App{
		Name:                 "parcel",
		Version:              Version,
		EnableBashCompletion: true,
	}
	app.Commands = []*cli.Command{
		cmd.NewPullCommand(cfg, logger),
		cmd.NewRenderCommand(cfg, logger),
	}

	app.Before = func(c *cli.Context) error {
		if c.Args().First() == "" {
			return nil
		}
		log.Infof("Running parcel version: %s", Version)
		return nil
	}
	app.Run(os.Args)
}
