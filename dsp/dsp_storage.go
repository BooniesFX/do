package dsp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/daseinio/do/common/log"
)

const fileNumberLimit = 1000

type fsNode struct {
	Name     string
	Size     int64
	Modified time.Time
	Children []*fsNode
}

func (s *Server) listFiles() *fsNode {
	rootDir := s.State.Config.DownloadDirectory
	root := &fsNode{}
	if info, err := os.Stat(rootDir); err == nil {
		if err := list(rootDir, info, root, new(int)); err != nil {
			log.Infof("File listing failed: %s", err)
		}
	}
	return root
}

//custom directory walk

func list(path string, info os.FileInfo, node *fsNode, n *int) error {
	if (!info.IsDir() && !info.Mode().IsRegular()) || strings.HasPrefix(info.Name(), ".") {
		return errors.New("Non-regular file")
	}
	(*n)++
	if (*n) > fileNumberLimit {
		return errors.New("Over file limit") //limit number of files walked
	}
	node.Name = info.Name()
	node.Size = info.Size()
	node.Modified = info.ModTime()
	if !info.IsDir() {
		return nil
	}
	children, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("Failed to list files")
	}
	node.Size = 0
	for _, i := range children {
		c := &fsNode{}
		p := filepath.Join(path, i.Name())
		if err := list(p, i, c, n); err != nil {
			continue
		}
		node.Size += c.Size
		node.Children = append(node.Children, c)
	}
	return nil
}
