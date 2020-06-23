package main

import (
	"fmt"
	"os"

	wzd_runner "github.com/infra-whizz/wzd/runner"

	wzlib_utils "github.com/infra-whizz/wzlib/utils"

	"github.com/infra-whizz/wzd"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

func runDummy(ctx *cli.Context) error {
	cli.ShowAppHelpAndExit(ctx, wzlib_utils.EX_USAGE)
	return nil
}

func runLocal(ctx *cli.Context) error {
	stateFile := ctx.String("state")

	if stateFile == "" {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			panic("This should not happen")
		}
		fmt.Println("Error: State file was not specified\n")
		os.Exit(wzlib_utils.EX_USAGE)
	}

	cms := wzd_runner.NewWzCMS()
	cms.SetStateRoot(ctx.String("state-dir"))
	cms.Call(stateFile)

	return nil
}

func runDaemon(ctx *cli.Context) error {
	daemon := wzd.NewWzDaemon()

	config := nanoconf.NewConfig(ctx.String("config"))
	cfgDaemon := config.Find("daemon")
	if cfgDaemon == nil {
		os.Exit(wzlib_utils.EX_USAGE)
		daemon.GetLogger().Errorf("Configuration for daemon was not found in file", ctx.String("config"))
	}

	daemon.SetPkiDirectory(config.Find("daemon").String("pki", ""))
	daemon.SetTraitsFile(config.Find("daemon").String("traits", ""))
	daemon.GetTransport().AddNatsServerURL(
		config.Find("transport").String("host", ""),
		config.Find("transport").DefaultInt("port", "", 4222))
	daemon.SetClusterFingerprint(config.Find("daemon").String("cluster-fingerprint", ""))
	daemon.SetupMachineIdUtil("")
	daemon.Run().AppLoop()

	return nil
}

func main() {
	appname := "wzd"
	confpath := nanoconf.NewNanoconfFinder(appname).DefaultSetup(nil)
	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "Whizz Client Daemon",
		Action:  runDummy,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to configuration file",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "state",
				Aliases:  []string{"s"},
				Usage:    "The name of the state",
				Required: false,
			},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "local",
			Usage:  "Run local state (in-place orchestration)",
			Action: runLocal,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "state-dir",
					Usage:   "Path to the static local root of states",
					Aliases: []string{"s"},
				},
			},
		},
		{
			Name:   "daemon",
			Usage:  "Run cluster daemon client",
			Action: runDaemon,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Path to configuration file",
					Required: false,
					Value:    confpath.SetDefaultConfig(confpath.FindFirst()).FindDefault(),
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
