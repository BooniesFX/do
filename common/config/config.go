package config

import (
	"time"
)

//VERSION set this via ldflags
var VERSION = "0"

const (
	NETWORK_ID_MAIN = 1
	NETWORK_ID_TEST = 2

	NETWORK_NAME_MAIN = "dasein"
	NETWORK_NAME_TEST = "test"
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
	DEFAULT_NETWORK_ADDRESS      = "tcp://localhost:6921"

	//component
	DEFAULT_BACKOFF_DELAY         = 5 * time.Second
	DEFAULT_BACKOFF_ATTEMPTS      = 5
	DEFAULT_BACKOFF_PRIORITY      = 100
	DEFAULT_SIGNATURE_POLICY      = "ed25519"
	DEFAULT_HASH_POLICY           = "blake2b"
	NETWORK_MAGIC_MAIN            = 0x19fd3b3a
	NETWORK_MAGIC_TEST            = 0x1aacb7f5
	DEFAULT_MAX_CONN_LIMIT        = 256
	DEFAULT_MAX_INBOUND_SINGLE_IP = 32
)

//default common parameter
const (
	DEFAULT_INIT_DIR     = ".do"
	DEFAULT_DATA_DIR     = "./data"
	DEFAULT_LOG_DIR      = "./log"
	DEFAULT_LOG_LEVEL    = 1  //INFO
	DEFAULT_MAX_LOG_SIZE = 20 //MB
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

type NetworkConfig struct {
	SeedList  []string
	NetworkId uint   `json:"id"`
	Magic     uint   `json:"n"`
	Name      string `help:"network name"`
}

//network config
type P2PConfig struct {
	Port           uint   `port`
	GRPCPort       uint   `help:"gRPC port to listen to"`
	JSONPort       uint   `help:"JSON-RPC port to listen to"`
	protocol       string `help:"protocol to use (kcp/tcp)"`
	Nat            bool   `help:"enable nat traversal"`
	DHT            bool   `help:"enable peer discovery"`
	Reconnect      bool   `help:"enable reconnections"`
	MaxConnLimit   uint
	MaxForSingleIP uint
	SignatureAlgo  string `help:"signature policy(ed25519)"`
	HashAlgo       string `help:"hash policy(blake2b)"`
}

type CommonConfig struct {
	LogLevel uint
	DataDir  string
}

type GenesisConfig struct {
}

type ConsensusConfig struct {
}

type DaseinConfig struct {
	Genesis   *GenesisConfig
	Common    *CommonConfig
	Consensus *ConsensusConfig
	Network   *NetworkConfig
	P2P       *P2PConfig
}

var MainNetWork = &NetworkConfig{
	SeedList:  MainNetSeeds,
	NetworkId: NETWORK_ID_MAIN,
	Magic:     NETWORK_MAGIC_MAIN,
	Name:      NETWORK_NAME_MAIN,
}

var TestNetWork = &NetworkConfig{
	SeedList:  TestNetSeeds,
	NetworkId: NETWORK_ID_TEST,
	Magic:     NETWORK_MAGIC_TEST,
	Name:      NETWORK_NAME_TEST,
}

func NewDaseinConfig(id uint) *DaseinConfig {
	net := &NetworkConfig{}
	switch id {
	case NETWORK_ID_TEST:
		net = TestNetWork
	case NETWORK_MAGIC_MAIN:
		net = MainNetWork
	default:

	}
	return &DaseinConfig{
		Genesis: &GenesisConfig{},
		Common: &CommonConfig{
			LogLevel: DEFAULT_LOG_LEVEL,
			DataDir:  DEFAULT_DATA_DIR,
		},
		Consensus: &ConsensusConfig{},
		Network:   net,
		P2P: &P2PConfig{
			Port:           DEFAULT_LISTEN_PORT,
			GRPCPort:       DEFAULT_GRPC_PORT,
			JSONPort:       DEFAULT_JSONRPC_PORT,
			protocol:       DEFAULT_LISTEN_PROTOCOL,
			Nat:            false,
			DHT:            true,
			Reconnect:      true,
			MaxConnLimit:   DEFAULT_MAX_CONN_LIMIT,
			MaxForSingleIP: DEFAULT_MAX_INBOUND_SINGLE_IP,
			SignatureAlgo:  DEFAULT_SIGNATURE_POLICY,
			HashAlgo:       DEFAULT_HASH_POLICY,
		},
	}
}

//current default config
var DefConfig = NewDaseinConfig(NETWORK_ID_TEST)
