package netserver

import (
	"github.com/daseinio/dasein-p2p/crypto/ed25519"
	p2p "github.com/daseinio/dasein-p2p/network"
	"github.com/daseinio/dasein-p2p/network/backoff"
	"github.com/daseinio/dasein-p2p/network/discovery"
	"github.com/daseinio/dasein-p2p/network/nat"
	"github.com/daseinio/dasein-p2p/peer"
	"github.com/daseinio/do/common/config"
	"github.com/daseinio/do/common/log"
	//msg "github.com/daseinio/do/netserver/p2pmsg"
	//jrpc "github.com/daseinio/do/netserver/rpc/json"
)

var netserver struct {
	peer *p2p.Network
	ID   *peer.ID
}

func (this *netserver) StartUp() {
	go this.peer.Listen()
	log.Infof("start listen at %s\n", this.peer.Address)
	this.peer.BlockUntilListening()
	seeds := config.DefConfig.Genesis.Network.SeedList
	if len(seeds) > 0 {
		this.peer.Bootstrap(seeds...)
	}
}

func (this *netserver) ShutDown() {
	this.peer.Close()
	log.Infoln("shuts down the entire network")
}

func (this *netserver) Initialize(magic uint32) {

	builder := p2p.NewBuilderWithOptions(
		p2p.SignaturePolicy(config.DEFAULT_SIGNATURE_POLICY),
		p2p.HashPolicy(config.DEFAULT_HASH_POLICY),
		p2p.ConnectionTimeout(config.DEFAULT_CONNECT_TIMEOUT),
		p2p.RecvWindowSize(config.DEFAULT_RECEIVE_WINDOWS_SIZE),
		p2p.SendWindowSize(config.DEFAULT_SEND_WINDOWS_SIZE),
		p2p.ReceiveBufferSize(config.DEFAULT_RECEIVE_BUFFER_SIZE),
		p2p.WriteBufferSize(config.DEFAULT_WRITE_SIZE),
		p2p.WriteFlushLatency(config.DEFAULT_WRITE_FLUSH_LATENCY),
		p2p.WriteTimeout(config.DEFAULT_WRITE_TIMEOUT),
	)
	keys := ed25519.RandomKeyPair()
	log.Infof("network public key: %s", keys.PublicKeyHex())
	builder.SetKeys(keys)
	protocol := config.DefConfig.P2P.Protocol
	host := "127.0.0.1"
	port := config.DefConfig.P2P.Port
	builder.SetAddress(p2p.FormatAddress(protocol, host, port))
	//component
	if config.DefConfig.P2P.Nat {
		nat.RegisterComponent(builder)
	}
	if config.DefConfig.P2P.DHT {
		builder.AddComponent(new(discovery.Component))
	}
	if config.DefConfig.P2P.Reconnect {
		builder.AddComponent(new(backoff.New(
			backoff.WithInitialDelay(config.DEFAULT_BACKOFF_DELAY),
			backoff.WithMaxAttempts(config.DEFAULT_BACKOFF_ATTEMPTS))))
	}
	//msg register
	//rpc setup
	//traffic register
	//
	//builder.SetMagic(magic)
	log.Infof("network magic number is %x", magic)
	this.peer, err = builder.Build()
	if err != nil {
		log.Fatalln(err)
		return
	}
	this.ID = this.peer.ID
	log.Infof("peer id is %x", this.ID)

}

func (this *netserver) Connect(addr string) bool {
	return true
}

func (this *netserver) Connect(id peer.ID) bool {
	return true
}

func (this *netserver) GetConnectPeers() {
	return this.net.ID
}

func (this *netserver) GetListenAddress() string {
	return this.net.Address
}

func (this *netserver) GetPubKey() *crypto.KeyPair {
	return this.net.GetKeys()
}

func (this *netserver) ConnectionsExist(addr string) bool {
	return this.net.ConnectionStateExists(addr)
}

//SendTo transfer msg to dest id, if connection not exsited, send to closest peers to relay it
func (this *netserver) SendTo(id *peer.ID, message proto.Message) error {
	// Check if we are the target.
	if this.ID.Equals(id) {
		return nil
	}
	//peer exsited
	if this.ConnectionsExist(id.Address) {
		this.peer.BroadcastByAddresses(message, id.Address)
		return nil
	}
	log.Infof("peer %s not connected, tranfer message by discovery router", id.String())
	//not in connection list
	Component, registered := this.peer.Component(discovery.ComponentID)
	if !registered {
		return errors.New("can`t reach destination id, discovery component not registered")
	}

	routes := Component.(*discovery.Component).Routes

	// Find the 3 closest peers from a nodes point of view (might include us).
	closestPeers := routes.FindClosestPeers(targetID, 3)

	// Remove self from the list.
	for i, id := range closestPeers {
		if id.Equals(this.ID) {
			closestPeers = append(closestPeers[:i], closestPeers[i+1:]...)
			break
		}
	}

	// Seems we have ran out of peers to attempt to propagate to.
	if len(closestPeers) == 0 {
		return errors.Errorf("could not found route from peer %d to peer %d", this.ID, id)
	}

	// Propagate message to the closest peer.
	this.peer.BroadcastByAddresses(message, closestPeers)
	return nil
}

//Broadcast send msg to a set of peer client by id
func (this *netserver) Broadcast(message proto.Message, ids ...peer.ID) {
	for _, id := range ids {
		err := this.SendTo(id, message)
		if err != nil {
			log.Errorln(err)
		}
	}
}
