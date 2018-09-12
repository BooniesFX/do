package dsp

import (
	"fmt"
	"path/filepath"

	"sync"
	"time"

	"github.com/daseinio/do/common"
	"github.com/daseinio/do/common/config"
	"github.com/daseinio/do/dsp/engine"
)

type Info struct {
	sync.Mutex
	Config    config.EngineConfig
	Downloads *fsNode
	Torrents  map[string]*engine.Torrent
	System    common.Stats
}

//Server is the "State" portion of the diagram
type Server struct {
	//torrent engine
	Engine *engine.DoEngine
	State  *Info
}

// Run the server
func (s *Server) Run() error {
	//torrent engine
	s.Engine = engine.New()
	//configure engine
	c := config.DefaultConfig.Engine

	if err := s.reconfigure(c); err != nil {
		return fmt.Errorf("initial configure failed: %s", err)
	}
	//poll torrents and files
	go func() {
		for {
			s.State.Lock()
			s.State.Torrents = s.Engine.GetTorrents()
			s.State.Downloads = s.listFiles()
			s.State.System.LoadStats(c.DownloadDirectory)
			s.State.Unlock()
			PrintState(s.State)
			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}

func (s *Server) reconfigure(c *config.EngineConfig) error {
	dldir, err := filepath.Abs(c.DownloadDirectory)
	if err != nil {
		return fmt.Errorf("Invalid path")
	}
	c.DownloadDirectory = dldir
	if err := s.Engine.Configure(c); err != nil {
		return err
	}
	return nil
}

func PrintState(obj *Info) {
	fmt.Println(obj)
}
