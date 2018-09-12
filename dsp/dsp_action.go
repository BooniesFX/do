package dsp

import (
	"bytes"

	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/daseinio/do/common/log"
)

func (s *Server) StartRemoteTorrent(url, name string) error {
	//convert url into torrent bytes
	remote, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Invalid remote torrent URL: %s (%s)", err, url)
	}
	data, err := ioutil.ReadAll(remote.Body)
	if err != nil {
		return fmt.Errorf("Failed to download remote torrent: %s", err)
	}
	filename := filepath.Join(s.State.Config.DownloadDirectory, name+".torrent")
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return fmt.Errorf("Failed to save torrent: %s", err)
	}
	log.Infof("torrent file downloaded:%s", filename)
	reader := bytes.NewBuffer(data)
	info, err := metainfo.Load(reader)
	if err != nil {
		return err
	}
	spec := torrent.TorrentSpecFromMetaInfo(info)
	if err := s.Engine.NewTorrent(spec); err != nil {
		return fmt.Errorf("Torrent error: %s", err)
	}
	return nil
}

func (s *Server) StartMagnet(url string) error {
	if err := s.Engine.NewMagnet(url); err != nil {
		return fmt.Errorf("Magnet error: %s", err)
	}
	return nil
}

func (s *Server) StartTask(infohash string) error {
	if err := s.Engine.StartTorrent(infohash); err != nil {
		return err
	}
	return nil
}

func (s *Server) DeleteTask(infohash string) error {
	if err := s.Engine.DeleteTorrent(infohash); err != nil {
		return err
	}
	return nil
}

func (s *Server) StartFromFile(name string) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Errorf("Read torrent file error: %s", err)
	}
	reader := bytes.NewBuffer(data)
	info, err := metainfo.Load(reader)
	if err != nil {
		return err
	}
	spec := torrent.TorrentSpecFromMetaInfo(info)
	if err := s.Engine.NewTorrent(spec); err != nil {
		return fmt.Errorf("Torrent error: %s", err)
	}
	return nil
}
