package main

import (
	"github.com/infra-whizz/wzd"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
	"os"
)

func run(ctx *cli.Context) error {
	config := nanoconf.NewConfig(ctx.String("config"))
	daemon := wzd.NewWzDaemon()
	daemon.GetTransport().AddNatsServerURL(
		config.Find("transport").String("host", ""),
		config.Find("transport").DefaultInt("port", "", 4222))
	daemon.Run()

	//cli.ShowAppHelpAndExit(ctx, 2)
	return nil
}

func main() {
	appname := "wzd"
	confpath := nanoconf.NewNanoconfFinder(appname).DefaultSetup(nil)
	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "Whizz Client Daemon",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to configuration file",
				Required: false,
				Value:    confpath.SetDefaultConfig(confpath.FindFirst()).FindDefault(),
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
