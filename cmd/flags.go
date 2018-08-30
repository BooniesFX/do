package cmd

import (
	"strings"

	"github.com/daseinio/do/common/config"
	"github.com/urfave/cli"
)

var (
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Use `<filename>` to specifies the config file. If doesn't specifies the config, do will use default config.",
	}
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~4). 0:stdout only 1:INFO 2:WARNING 3:ERROR 4:FATAL",
		Value: config.DEFAULT_LOG_LEVEL,
	}
	DataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Using dir `<path>` to storage block data",
		Value: config.DEFAULT_DATA_DIR,
	}
	//P2P setting
	NetworkIdFlag = cli.UintFlag{
		Name:  "networkid",
		Usage: "Using to specify the network ID. Different networkids cannot connect to the network. 1:dasein main net, 2: test net, and other for custom network",
		Value: config.NETWORK_ID_MAIN,
	}
	NetworkNameFlag = cli.StringFlag{
		Name:  "networkname",
		Usage: "Using to specify the network name for custom network",
		Value: config.NETWORK_NAME_TEST,
	}
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
	GRPCPortFlag = cli.UintFlag{
		Name:  "grpcport",
		Usage: "Using to specify the gRPC port number",
		Value: config.DEFAULT_GRPC_PORT,
	}
	JSONPortFlag = cli.UintFlag{
		Name:  "jsonport",
		Usage: "Using to specify the JSON-RPC network port number",
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
	MaxConnectionFlag = cli.UintFlag{
		Name:  "maxconnect",
		Usage: "Max connection in total",
		Value: config.DEFAULT_MAX_CONN_LIMIT,
	}
	MaxForSingleIPFlag = cli.UintFlag{
		Name:  "maxinforsingleip",
		Usage: "Max connection in bound for single ip",
		Value: config.DEFAULT_MAX_INBOUND_SINGLE_IP,
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
