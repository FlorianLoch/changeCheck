package persistence

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"path"

	"github.com/spf13/afero"
	"golang.org/x/net/html"
)

type Persistor interface {
	Store(url string, xpath string, nodes []*html.Node) error
	Load(url string, xpath string) ([]*html.Node, error)
}

const (
	cacheDirName = ".change_check_cache"
)

type FSPersistor struct {
	fs afero.Fs
}

func NewFSPersistor(fs afero.Fs) *FSPersistor {
	return &FSPersistor{
		fs: fs,
	}
}

func (f *FSPersistor) Store(url string, xpath string, nodes []*html.Node) error {
	if len(nodes) == 0 {
		return nil
	}

	cacheDir, fileName, err := paths(url, xpath)
	if err != nil {
		return err
	}

	err = f.fs.MkdirAll(cacheDir, 755)
	if err != nil {
		return err
	}

	// Create file, if it exists overwrite its content
	file, err := f.fs.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(nodes)
	if err != nil {
		return err
	}

	return file.Close()
}

func (f *FSPersistor) Load(url string, xpath string) ([]*html.Node, error) {
	_, fileName, err := paths(url, xpath)
	if err != nil {
		return nil, err
	}

	file, err := f.fs.Open(fileName)
	if err != nil {
		if err == afero.ErrFileNotFound {
			return make([]*html.Node, 0), nil
		}

		return nil, err
	}

	var nodes []*html.Node
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func paths(url, xpath string) (string, string, error) {
	cacheDir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	cacheDir = path.Join(cacheDir, cacheDirName)

	fileName := fmt.Sprintf("%x", sha256.Sum256([]byte(url+xpath)))
	fileName = path.Join(cacheDir, fileName)

	return cacheDir, fileName, nil
}
