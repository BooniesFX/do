package config

import (
	"crypto/sha256"
	"encoding/binary"
	"time"
)

//VERSION set this via ldflags
var VERSION = ""

const (
	NETWORK_NAME_MAIN = "dasein"
	NETWORK_NAME_TEST = "test"
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
	DEFAULT_WRITE_SIZE           = 16 * 1024
	DEFAULT_RECEIVE_BUFFER_SIZE  = 16 * 1024 * 1024
	DEFAULT_WRITE_FLUSH_LATENCY  = 50 * time.Millisecond
	DEFAULT_WRITE_TIMEOUT        = 3 * time.Second
	DEFAULT_LISTEN_PORT          = uint(6921)
	DEFAULT_GRPC_PORT            = uint(6922)
	DEFAULT_JSONRPC_PORT         = uint(8923)
	DEFAULT_LISTEN_PROTOCOL      = "tcp"
	DEFAULT_SIGNATURE_POLICY     = "ed25519"
	DEFAULT_HASH_POLICY          = "blake2b"

	//component
	DEFAULT_BACKOFF_DELAY    = 5 * time.Second
	DEFAULT_BACKOFF_ATTEMPTS = 5

	DEFAULT_MAX_CONN_INBOUND_LIMIT  = 128
	DEFAULT_MAX_CONN_OUTBOUND_LIMIT = 32
	DEFAULT_MAX_INBOUND_SINGLE_IP   = 16
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
	Name     string   `json:"network name"`
}

//network config
type P2PConfig struct {
	Port             uint
	Protocol         string
	Nat              bool
	DHT              bool
	Reconnect        bool
	MaxConnInLimit   uint
	MaxConnOutLimit  uint
	MaxInForSingleIP uint
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

//return hash of genesis config
func (this *GenesisConfig) genesisMagic() uint32 {
	//generate genesis byte
	//TODO code if compute from genesis file
	code := []byte{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2}
	data := sha256.Sum256(code)
	return binary.LittleEndian.Uint32(data[0:4])
}

type RpcConfig struct {
	EnableGRPC bool
	EnableJson bool
	GRPCPort   uint
	JSONPort   uint
}

type EngineConfig struct {
	IncomingPort      uint
	DownloadDirectory string
}

type DaseinConfig struct {
	Magic     uint32
	Genesis   *GenesisConfig
	Consensus *ConsensusConfig
	Common    *CommonConfig
	P2P       *P2PConfig
	Rpc       *RpcConfig
	Engine    *EngineConfig
}

var MainNetGenesis = &GenesisConfig{
	Network: &SeedNetworkConfig{
		SeedList: MainNetSeeds,
		Name:     NETWORK_NAME_MAIN,
	},
}

var TestNetGenesis = &GenesisConfig{
	Network: &SeedNetworkConfig{
		SeedList: TestNetSeeds,
		Name:     NETWORK_NAME_TEST,
	},
}

func DefConfig() *DaseinConfig {
	return &DaseinConfig{
		Common: &CommonConfig{
			LogLevel:  DEFAULT_LOG_LEVEL,
			LogStderr: false,
		},
		Consensus: &ConsensusConfig{},
		P2P: &P2PConfig{
			Port:             DEFAULT_LISTEN_PORT,
			Protocol:         DEFAULT_LISTEN_PROTOCOL,
			Nat:              true,
			DHT:              true,
			Reconnect:        true,
			MaxConnInLimit:   DEFAULT_MAX_CONN_INBOUND_LIMIT,
			MaxConnOutLimit:  DEFAULT_MAX_CONN_OUTBOUND_LIMIT,
			MaxInForSingleIP: DEFAULT_MAX_INBOUND_SINGLE_IP,
		},
		Rpc: &RpcConfig{
			EnableGRPC: false,
			EnableJson: false,
			GRPCPort:   DEFAULT_GRPC_PORT,
			JSONPort:   DEFAULT_JSONRPC_PORT,
		},
		Engine: &EngineConfig{
			IncomingPort:      DEFAULT_ENGINE_INCOMINGPORT,
			DownloadDirectory: DEFAULT_ENGINE_DOWNLOADDIR,
		},
	}
}

//current default config
var DefaultConfig = DefConfig()

func (this *DaseinConfig) Set(n *GenesisConfig) {
	this.Genesis.Network = n.Network
	this.Magic = n.genesisMagic()
}
