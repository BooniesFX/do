package cmd

import (
	"strings"

	"github.com/daseinio/do/common/config"
	"github.com/urfave/cli"
)

var (
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Use `<filename>` to specifies the config file to connect to cunstomize network. If doesn't specifies the config, do will use default config(mainnet).",
	}
	TestNetFlag = cli.BoolFlag{
		Name:  "testnet",
		Usage: "use test net config , default use mainnet config",
	}
	//commmon
	LogStderrFlag = cli.BoolFlag{
		Name:  "logstderr",
		Usage: "log to standard error instead of files,default false",
	}
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~4). 0:DEBUG 1:INFO 2:WARNING 3:ERROR 4:FATAL",
		Value: config.DEFAULT_LOG_LEVEL,
	}
	//P2P setting
	ProtocolFlag = cli.StringFlag{
		Name:  "protocol",
		Usage: "Using to specify the protocol for listen(tcp/kcp)",
		Value: config.DEFAULT_LISTEN_PROTOCOL,
	}
	PortFlag = cli.UintFlag{
		Name:  "port",
		Usage: "Using to specify the P2P network port number",
		Value: config.DEFAULT_LISTEN_PORT,
	}
	NatSupportFlag = cli.BoolFlag{
		Name:  "nat",
		Usage: "enable nat. ",
	}
	DHTSupportFlag = cli.BoolFlag{
		Name:  "dht",
		Usage: "enable DHT. ",
	}
	BackOffSupportFlag = cli.BoolFlag{
		Name:  "reconnect",
		Usage: "reconnect enable the client try to reconnect remote peer whose connect is broken. ",
	}
	MaxInBoundConnectionFlag = cli.UintFlag{
		Name:  "maxinbound",
		Usage: "Max connection inbound",
		Value: config.DEFAULT_MAX_CONN_INBOUND_LIMIT,
	}
	MaxOutBoundConnectionFlag = cli.UintFlag{
		Name:  "maxoutbound",
		Usage: "Max connection outbound",
		Value: config.DEFAULT_MAX_CONN_OUTBOUND_LIMIT,
	}
	MaxForSingleIPFlag = cli.UintFlag{
		Name:  "maxinboundforsingleip",
		Usage: "Max connection in bound for single ip",
		Value: config.DEFAULT_MAX_INBOUND_SINGLE_IP,
	}
	//RPC setting
	EnableGRPCFlag = cli.BoolFlag{
		Name:  "grpc",
		Usage: "enable gRPC",
	}
	GRPCPortFlag = cli.UintFlag{
		Name:  "grpcport",
		Usage: "Using to specify the gRPC port number",
		Value: config.DEFAULT_GRPC_PORT,
	}
	EnableJsonFlag = cli.BoolFlag{
		Name:  "jsonrpc",
		Usage: "enable JSON-RPC",
	}
	JSONPortFlag = cli.UintFlag{
		Name:  "jsonport",
		Usage: "Using to specify the JSON-RPC network port number",
		Value: config.DEFAULT_JSONRPC_PORT,
	}
	//Engine setting
	DownloadDirFlag = cli.StringFlag{
		Name:  "dir",
		Usage: "Using dir `<path>` to storage downlod data",
		Value: config.DEFAULT_ENGINE_DOWNLOADDIR,
	}
	IncomigPortFlag = cli.UintFlag{
		Name:  "engineport",
		Usage: "Using to specify the engine port number",
		Value: config.DEFAULT_ENGINE_INCOMINGPORT,
	}
	//bt
	RemoteDownloadDirFlag = cli.StringFlag{
		Name:  "url",
		Usage: "Using url `<url_to_torrent>` to download file by remote torrent file",
	}
	DownloadFlag = cli.StringFlag{
		Name:  "open",
		Usage: "Using start `<path_to_torrent>` to download file by loacal torrent file",
	}
	MagnetFlag = cli.StringFlag{
		Name:  "magnet",
		Usage: "Using magnet `<url_to_magnet>` to download file by magnet",
	}
)

//GetFlagName deal with short flag, and return the flag name whether flag name have short name
func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
