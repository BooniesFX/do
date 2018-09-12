package engine

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/daseinio/do/common/config"
)

type DoEngine struct {
	mut       sync.Mutex
	client    *torrent.Client
	defconfig *config.EngineConfig
	ts        map[string]*Torrent
}

func New() *DoEngine {
	return &DoEngine{ts: map[string]*Torrent{}}
}

func (e *DoEngine) Config() *config.EngineConfig {
	return e.defconfig
}

func (e *DoEngine) Configure(c *config.EngineConfig) error {
	//recieve config
	if e.client != nil {
		e.client.Close()
		time.Sleep(1 * time.Second)
	}
	if c.IncomingPort <= 0 {
		return fmt.Errorf("Invalid incoming port (%d)", c.IncomingPort)
	}
	tc := torrent.ClientConfig{
		DhtStartingNodes: dht.GlobalBootstrapAddrs,
		DataDir:          c.DownloadDirectory,
		NoUpload:         !c.EnableUpload,
		Seed:             c.EnableSeeding,
	}
	tc.DisableEncryption = c.DisableEncryption
	tc.SetListenAddr(":53007")
	client, err := torrent.NewClient(&tc)
	if err != nil {
		return err
	}
	e.mut.Lock()
	e.defconfig = c
	e.client = client
	e.mut.Unlock()
	//reset
	e.GetTorrents()
	return nil
}

func (e *DoEngine) NewMagnet(magnetURI string) error {
	tt, err := e.client.AddMagnet(magnetURI)
	if err != nil {
		return err
	}
	return e.newTorrent(tt)
}

func (e *DoEngine) NewTorrent(spec *torrent.TorrentSpec) error {
	tt, _, err := e.client.AddTorrentSpec(spec)
	if err != nil {
		return err
	}
	return e.newTorrent(tt)
}

func (e *DoEngine) newTorrent(tt *torrent.Torrent) error {
	t := e.upsertTorrent(tt)
	go func() {
		<-t.t.GotInfo()
		e.StartTorrent(t.InfoHash)
	}()
	return nil
}

//GetTorrents moves torrents out of the anacrolix/torrent
//and into the local cache
func (e *DoEngine) GetTorrents() map[string]*Torrent {
	e.mut.Lock()
	defer e.mut.Unlock()

	if e.client == nil {
		return nil
	}
	for _, tt := range e.client.Torrents() {
		e.upsertTorrent(tt)
	}
	return e.ts
}

func (e *DoEngine) upsertTorrent(tt *torrent.Torrent) *Torrent {
	ih := tt.InfoHash().HexString()
	torrent, ok := e.ts[ih]
	if !ok {
		torrent = &Torrent{InfoHash: ih}
		e.ts[ih] = torrent
	}
	//update torrent fields using underlying torrent
	torrent.Update(tt)
	return torrent
}

func (e *DoEngine) getTorrent(infohash string) (*Torrent, error) {
	ih, err := str2ih(infohash)
	if err != nil {
		return nil, err
	}
	t, ok := e.ts[ih.HexString()]
	if !ok {
		return t, fmt.Errorf("Missing torrent %x", ih)
	}
	return t, nil
}

func (e *DoEngine) getOpenTorrent(infohash string) (*Torrent, error) {
	t, err := e.getTorrent(infohash)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (e *DoEngine) StartTorrent(infohash string) error {
	t, err := e.getOpenTorrent(infohash)
	if err != nil {
		return err
	}
	if t.Started {
		return fmt.Errorf("Already started")
	}
	t.Started = true
	for _, f := range t.Files {
		if f != nil {
			f.Started = true
		}
	}
	if t.t.Info() != nil {
		t.t.DownloadAll()
	}
	return nil
}

func (e *DoEngine) DeleteTorrent(infohash string) error {
	t, err := e.getTorrent(infohash)
	if err != nil {
		return err
	}
	//os.Remove(filepath.Join(e.cacheDir, infohash+".torrent"))
	delete(e.ts, t.InfoHash)
	ih, _ := str2ih(infohash)
	if tt, ok := e.client.Torrent(ih); ok {
		tt.Drop()
	}
	return nil
}

func str2ih(str string) (metainfo.Hash, error) {
	var ih metainfo.Hash
	e, err := hex.Decode(ih[:], []byte(str))
	if err != nil {
		return ih, fmt.Errorf("Invalid hex string")
	}
	if e != 20 {
		return ih, fmt.Errorf("Invalid length")
	}
	return ih, nil
}
