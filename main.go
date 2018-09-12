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
	"github.com/daseinio/do/dsp"
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
		cmd.TestNetFlag,
		cmd.LogStderrFlag,
		cmd.LogLevelFlag,
		//p2p setting
		cmd.ProtocolFlag,
		cmd.PortFlag,
		cmd.NatSupportFlag,
		cmd.DHTSupportFlag,
		cmd.BackOffSupportFlag,
		cmd.MaxInBoundConnectionFlag,
		cmd.MaxOutBoundConnectionFlag,
		cmd.MaxForSingleIPFlag,
		//rpc
		cmd.EnableGRPCFlag,
		cmd.GRPCPortFlag,
		cmd.EnableJsonFlag,
		cmd.JSONPortFlag,
		//engine
		cmd.DownloadDirFlag,
		cmd.IncomigPortFlag,
		//bt
		cmd.RemoteDownloadDirFlag,
		cmd.DownloadFlag,
		cmd.MagnetFlag,
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

	err := initConfig(ctx)
	if err != nil {
		log.Errorf("initConfig error:%s", err)
		return
	}
	svr := &dsp.Server{
		State: &dsp.Info{},
	}
	if err = svr.Run(); err != nil {
		log.Errorf("run bt server error:%s", err)
		return
	}
	url := ctx.String(cmd.GetFlagName(cmd.RemoteDownloadDirFlag))
	if url != "" {
		svr.StartRemoteTorrent(url, "test")
	}
	url = ctx.String(cmd.GetFlagName(cmd.MagnetFlag))
	log.Infoln(url)
	if url != "" {
		svr.StartMagnet(url)
	}
	waitToExit()
}

func initConfig(ctx *cli.Context) error {
	//init ontology config from cli
	_, err := cmd.SetDaseinConfig(ctx)
	if err != nil {
		return err
	}
	log.Infof("Config init success")
	return nil
}

func initLog(ctx *cli.Context) {
	//init log module
	log.SetLevel(ctx.GlobalUint(cmd.GetFlagName(cmd.LogLevelFlag)))
	log.SetMaxSize(config.DEFAULT_MAX_LOG_SIZE)
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
