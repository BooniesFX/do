package netserver

import (
	"github.com/daseinio/dasein-p2p/crypto/ed25519"
	p2p "github.com/daseinio/dasein-p2p/network"
	"github.com/daseinio/do/common/config"
	"github.com/daseinio/do/common/log"
	//msg "github.com/daseinio/do/netserver/p2pmsg"
	//jrpc "github.com/daseinio/do/netserver/rpc/json"
)

var netserver struct {
	net *p2p.Network
}

func (this *netserver) StartUp() {

}

func (this *netserver) ShutDown() {

}

func (this *netserver) Initialize() {
	this.net.Init()
	builder := NewBuilderWithOptions(
		SignaturePolicy(signaturePolicy),
		HashPolicy(hashPolicy),
	)
	keys := ed25519.RandomKeyPair()
	log.Infof("network public key: %s", keys.PublicKeyHex())
	builder.SetKeys(keys)
	protocol := config.DefConfig.P2P.Protocol
	host := "127.0.0.1"
	port := config.DefConfig.P2P.
		builder.SetAddress(this.net.FormatAddress(protocol, host, port))
}

func (this *netserver) Connect() {

}

func (this *netserver) GetNetID() {
	return this.net.ID
}

func (this *netserver) GetListenAddress() string {
	return this.net.Address
}

func (this *netserver) GetPubKey() *crypto.KeyPair {
	return this.net.GetKeys()
}

func (this *netserver) ConnectionExist(addr string) bool {
	return this.net.ConnectionStateExists(addr)
}
