package cmd

import (
	"github.com/daseinio/do/common"
	"github.com/daseinio/do/common/config"
	"github.com/daseinio/do/common/log"
	"github.com/urfave/cli"
)

func SetDaseinConfig(ctx *cli.Context) (*config.DaseinConfig, error) {
	cfg := config.DefConfig()
	err := setGenesis(ctx, cfg)
	if err != nil {
		return nil, err
	}
	setCommonConfig(ctx, cfg.Common)
	setConsensusConfig(ctx, cfg.Consensus)
	setP2PNodeConfig(ctx, cfg.P2P)
	setRpcConfig(ctx, cfg.Rpc)
	setEngineConfig(ctx, cfg.Engine)
	return cfg, nil
}

func setGenesis(ctx *cli.Context, cfg *config.DaseinConfig) error {
	testnet := ctx.Bool(GetFlagName(TestNetFlag))
	if testnet {
		cfg.Set(config.TestNetGenesis)
	} else {
		cfg.Set(config.MainNetGenesis)
	}

	if !ctx.IsSet(GetFlagName(ConfigFlag)) {
		return nil
	}

	genesisFile := ctx.String(GetFlagName(ConfigFlag))
	if !common.FileExisted(genesisFile) {
		return nil
	}

	newGenesisCfg := &config.GenesisConfig{}
	err := common.GetJsonObjectFromFile(genesisFile, newGenesisCfg)
	if err != nil {
		return err
	}
	log.Infof("Load genesis config:%s", genesisFile)
	cfg.Set(newGenesisCfg)

	return nil
}

func setCommonConfig(ctx *cli.Context, cfg *config.CommonConfig) {
	cfg.LogLevel = ctx.Uint(GetFlagName(LogLevelFlag))
	cfg.LogStderr = ctx.Bool(GetFlagName(LogStderrFlag))

}

func setConsensusConfig(ctx *cli.Context, cfg *config.ConsensusConfig) {
	//
}

func setP2PNodeConfig(ctx *cli.Context, cfg *config.P2PConfig) {

	cfg.Port = ctx.Uint(GetFlagName(PortFlag))
	cfg.Protocol = ctx.String(GetFlagName(ProtocolFlag))
	cfg.Nat = ctx.Bool(GetFlagName(NatSupportFlag))
	cfg.DHT = ctx.Bool(GetFlagName(DHTSupportFlag))
	cfg.Reconnect = ctx.Bool(GetFlagName(BackOffSupportFlag))
	cfg.MaxConnInLimit = ctx.Uint(GetFlagName(MaxInBoundConnectionFlag))
	cfg.MaxConnOutLimit = ctx.Uint(GetFlagName(MaxOutBoundConnectionFlag))
	cfg.MaxInForSingleIP = ctx.Uint(GetFlagName(MaxForSingleIPFlag))

}

func setRpcConfig(ctx *cli.Context, cfg *config.RpcConfig) {
	cfg.EnableGRPC = ctx.Bool(GetFlagName(EnableGRPCFlag))
	cfg.GRPCPort = ctx.Uint(GetFlagName(GRPCPortFlag))
	cfg.EnableJson = ctx.Bool(GetFlagName(EnableJsonFlag))
	cfg.JSONPort = ctx.Uint(GetFlagName(JSONPortFlag))
}

func setEngineConfig(ctx *cli.Context, cfg *config.EngineConfig) {
	cfg.IncomingPort = ctx.Uint(GetFlagName(IncomigPortFlag))
	cfg.DownloadDirectory = ctx.String(GetFlagName(DownloadDirFlag))
}
