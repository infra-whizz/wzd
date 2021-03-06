package main

import (
	"fmt"
	"os"
	"strings"

	wzd_runner "github.com/infra-whizz/wzd/runner"
	"github.com/sirupsen/logrus"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"

	"github.com/infra-whizz/wzd"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

func setLogger(ctx *cli.Context) {
	if ctx.String("logformat") == "json" {
		wzlib_logger.SetCurrentLogger(wzlib_logger.GetJSONLogger(logrus.DebugLevel, nil))
	}
}

func runDummy(ctx *cli.Context) error {
	cli.ShowAppHelpAndExit(ctx, wzlib_utils.EX_USAGE)
	return nil
}

func containerProxy() {
	// :c-<id>:chroot:module:jsonpath
	if len(os.Args) > 2 && strings.HasPrefix(os.Args[2], ":c") {
		args := strings.Split(os.Args[2], ":")
		if len(args) == 5 {
			args = args[1:]
			conf := new(wzlib_utils.WzContainerParam)
			conf.Root = args[1]
			conf.Command = args[2]
			conf.Args = []string{args[3]}
			stdout, stderr, err := wzlib_utils.NewWzContainer(conf).Container()

			os.Stdout.WriteString(stdout + "\n")
			os.Stderr.WriteString(stderr + "\n")
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("Error running container: %s", err.Error()))
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			os.Stderr.WriteString(fmt.Sprintf("Error parsing internal container args: %d", len(args)))
			os.Exit(1)
		}
	}
}

func runLocal(ctx *cli.Context) error {
	setLogger(ctx)

	stateId := ctx.String("state")
	stateDir := ctx.String("dir")

	containerProxy()

	if stateId == "" || stateDir == "" {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			panic("This should not happen")
		}
		var what string
		if stateId == "" {
			what = "file"
		} else {
			what = "root directory"
		}
		fmt.Printf("Error: State %s was not specified", what)
		os.Exit(wzlib_utils.EX_USAGE)
	}

	config := nanoconf.NewConfig(ctx.String("config"))
	cms := wzd_runner.NewWzCMS(stateDir).
		SetPyInterpreter(config.Find("ansible").String("python", "")).
		SetChrootedModules(ctx.String("modules-root"))
	_, res, _ := cms.OfflineCallById(stateId)
	for _, logEntry := range res {
		logEntry.Log()
	}

	return nil
}

func runDaemon(ctx *cli.Context) error {
	setLogger(ctx)
	daemon := wzd.NewWzDaemon()

	config := nanoconf.NewConfig(ctx.String("config"))
	cfgDaemon := config.Find("daemon")
	if cfgDaemon == nil {
		os.Exit(wzlib_utils.EX_USAGE)
		daemon.GetLogger().Errorf("Configuration for daemon was not found in file: %s", ctx.String("config"))
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
				Name:    "logformat",
				Aliases: []string{"f"},
				Usage:   "Log format. Choices: default, json",
				Value:   "default",
			},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "local",
			Usage:  "Run local state (in-place orchestration)",
			Action: runLocal,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dir",
					Usage:   "Path to the static local root of states",
					Aliases: []string{"d"},
				},
				&cli.StringFlag{
					Name:    "state",
					Usage:   "The name of the state",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:    "modules-root",
					Usage:   "Run all modules within an alternative root",
					Value:   "/",
					Aliases: []string{"r"},
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
