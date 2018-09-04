package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/daseinio/do/cmd"
	"github.com/daseinio/do/common/config"
	"github.com/daseinio/do/common/log"
	"github.com/urfave/cli"
)

func initAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "dasein client"
	app.Action = doit
	app.Version = config.VERSION
	app.Copyright = "Copyright in 2018 The Dasein Authors"
	app.Commands = []cli.Command{}
	app.Flags = []cli.Flag{
		//common setting
		cmd.ConfigFlag,
		cmd.LogStderrFlag,
		cmd.LogLevelFlag,
		cmd.DataDirFlag,
		//p2p setting
		cmd.ProtocolFlag,
		cmd.PortFlag,
		cmd.GRPCPortFlag,
		cmd.JSONPortFlag,
		cmd.NatSupportFlag,
		cmd.DHTSupportFlag,
		cmd.BackOffSupportFlag,
		cmd.MaxConnectionFlag,
		cmd.MaxForSingleIPFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := initAPP().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func doit(ctx *cli.Context) {
	initLog(ctx)
	waitToExit()
}

func initLog(ctx *cli.Context) {
	//init log module
	log.SetLevel(ctx.GlobalUint(cmd.GetFlagName(cmd.LogLevelFlag)))
	log.SetMaxSize(DEFAULT_MAX_LOG_SIZE)
	if ctx.Bool(cmd.GetFlagName(cmd.LogStderrFlag)) {
		log.InitLog(0, config.DEFAULT_LOG_DIR)
	} else {
		log.InitLog(1, config.DEFAULT_LOG_DIR)
	}
	log.Info("start logging...")
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("do received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
