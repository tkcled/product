package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	config "hshelby-tkcled-product/config"
	"hshelby-tkcled-product/src/database"
	"hshelby-tkcled-product/src/server"
)

var (
	configPrefix string
	configSource string
)

func main() {
	app := cli.NewApp()
	app.Name = "Product microservice"
	app.Usage = "Product microservice"
	app.Copyright = "Copyright © 2024 HShelby Groups. All Rights Reserved."
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "configPrefix",
			Aliases:     []string{"confPrefix"},
			Usage:       "prefix for config",
			Value:       "product",
			Destination: &configPrefix,
		},
		&cli.StringFlag{
			Name:        "configSource",
			Aliases:     []string{"confSource"},
			Value:       "../config/.env",
			Usage:       "set path to environment file",
			Destination: &configSource,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "Start the product server",
			Action: Serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr-graph",
					Aliases: []string{"address-graph"},
					Value:   "0.0.0.0:8081",
					Usage:   "address for serve graph",
				},
				&cli.StringFlag{
					Name:    "addr-grpc",
					Aliases: []string{"address-grpc"},
					Value:   "0.0.0.0:9091",
					Usage:   "address for serve grpc",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		return config.LoadFromEnv(configPrefix, configSource)
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endSignal := make(chan os.Signal, 1)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func(ctx context.Context, errChan chan error) {
		err := app.RunContext(ctx, os.Args)
		errChan <- err
	}(ctx, errChan)

	select {
	case sign := <-endSignal:
		log.Println("shutting down. reason:", sign)
		return
	case err := <-errChan:
		if err == nil {
			return
		}
		log.Println("encountered error:", err)
		return
	}

}

func Serve(c *cli.Context) error {
	ctx := c.Context
	err := database.ConnectDatabse(ctx)
	if err != nil {
		panic(err)
	}

	return server.ServeGraph(c.Context, c.String("addr-graph"))
}
