package config

import (
	"time"
)

//VERSION set this via ldflags
var VERSION = ""

const (
	NETWORK_NAME_MAIN  = "dasein"
	NETWORK_NAME_TEST  = "test"
	NETWORK_MAGIC_MAIN = 0x19fd3b3a
	NETWORK_MAGIC_TEST = 0x1aacb7f5
)

type nettype int32

const (
	mainnet nettype = iota
	testnet
	custom
)

//default p2p parameter
const (
	DEFAULT_CONNECT_TIMEOUT      = 15 * time.Second
	DEFAULT_RECEIVE_WINDOWS_SIZE = 4096
	DEFAULT_SEND_WINDOWS_SIZE    = 4096
	DEFAULT_WRITE_FLUSH_LATENCY  = 50 * time.Millisecond
	DEFAULT_WRITE_TIMEOUT        = 3 * time.Second
	DEFAULT_LISTEN_PORT          = uint(6921)
	DEFAULT_GRPC_PORT            = uint(6922)
	DEFAULT_JSONRPC_PORT         = uint(8923)
	DEFAULT_LISTEN_PROTOCOL      = "tcp"

	//component
	DEFAULT_BACKOFF_DELAY         = 5 * time.Second
	DEFAULT_BACKOFF_ATTEMPTS      = 5
	DEFAULT_BACKOFF_PRIORITY      = 100
	DEFAULT_SIGNATURE_POLICY      = "ed25519"
	DEFAULT_HASH_POLICY           = "blake2b"
	DEFAULT_MAX_CONN_LIMIT        = 256
	DEFAULT_MAX_INBOUND_SINGLE_IP = 32
)

//default common parameter
const (
	DEFAULT_INIT_DIR     = ".do"
	DEFAULT_LOG_DIR      = "./log"
	DEFAULT_LOG_LEVEL    = 1                //INFO
	DEFAULT_MAX_LOG_SIZE = 20 * 1024 * 1024 //MB
)

//default engine parameter
const (
	DEFAULT_ENGINE_INCOMINGPORT = uint(16921)
	DEFAULT_ENGINE_DOWNLOADDIR  = "./data"
)

//main net genesis config
var MainNetSeeds = []string{
	"tcp://seed1.dasein.io:6921",
	"tcp://seed2.dasein.io:6921",
	"tcp://seed3.dasein.io:6921",
	"tcp://seed4.dasein.io:6921",
	"tcp://seed5.dasein.io:6921"}

var TestNetSeeds = []string{
	"kcp://seed1.dasein.io:7921",
	"kcp://seed2.dasein.io:7921",
	"kcp://seed3.dasein.io:7921",
	"kcp://seed4.dasein.io:7921",
	"kcp://seed5.dasein.io:7921"}

type SeedNetworkConfig struct {
	SeedList []string `json:"seeds"`
	Magic    uint     `json:"magic"`
	Name     string   `json:"network name"`
}

//network config
type P2PConfig struct {
	Port           uint
	Protocol       string
	Nat            bool
	DHT            bool
	Reconnect      bool
	MaxConnLimit   uint
	MaxForSingleIP uint
	//SignatureAlgo  string
	//HashAlgo       string
}

type CommonConfig struct {
	LogLevel  uint
	LogStderr bool
}

type GenesisConfig struct {
	Network *SeedNetworkConfig
}

type ConsensusConfig struct {
}

type RpcConfig struct {
	EnableGRPC bool
	EnableJson bool
	GRPCPort   uint
	JSONPort   uint
}

type EngineConfig struct {
	EnableUpload      bool
	EnableSeeding     bool
	IncomingPort      int
	DownloadDirectory string
}

type DaseinConfig struct {
	Genesis   *GenesisConfig
	Common    *CommonConfig
	Consensus *ConsensusConfig
	P2P       *P2PConfig
	Rpc       *RpcConfig
	Engine    *EngineConfig
}

var MainNetWork = &SeedNetworkConfig{
	SeedList: MainNetSeeds,
	Magic:    NETWORK_MAGIC_MAIN,
	Name:     NETWORK_NAME_MAIN,
}

var TestNetWork = &SeedNetworkConfig{
	SeedList: TestNetSeeds,
	Magic:    NETWORK_MAGIC_TEST,
	Name:     NETWORK_NAME_TEST,
}

func NewDaseinConfig(nt nettype) *DaseinConfig {
	net := &SeedNetworkConfig{}
	switch nt {
	case mainnet:
		net = MainNetWork
	case testnet:
		net = TestNetWork
	default: //custom
	}
	return &DaseinConfig{
		Genesis: &GenesisConfig{
			Network: net,
		},
		Common: &CommonConfig{
			LogLevel:  DEFAULT_LOG_LEVEL,
			LogStderr: false,
		},
		Consensus: &ConsensusConfig{},
		P2P: &P2PConfig{
			Port:           DEFAULT_LISTEN_PORT,
			Protocol:       DEFAULT_LISTEN_PROTOCOL,
			Nat:            true,
			DHT:            true,
			Reconnect:      true,
			MaxConnLimit:   DEFAULT_MAX_CONN_LIMIT,
			MaxForSingleIP: DEFAULT_MAX_INBOUND_SINGLE_IP,
			//SignatureAlgo:  DEFAULT_SIGNATURE_POLICY,
			//HashAlgo:       DEFAULT_HASH_POLICY,
		},
		Rpc: &RpcConfig{
			EnableGRPC: false,
			EnableJson: false,
			GRPCPort:   DEFAULT_GRPC_PORT,
			JSONPort:   DEFAULT_JSONRPC_PORT,
		},
		Engine: &EngineConfig{},
	}
}

//current default config
var DefConfig = NewDaseinConfig(testnet)
