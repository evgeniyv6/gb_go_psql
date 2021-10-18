package configuration

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type fileSystem interface {
	Open(name string) (ifile, error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
}

type ifile interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type OsFileSystem struct{}

func (OsFileSystem) Open(name string) (ifile, error) {
	return os.Open(name)
}
func (OsFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
func (OsFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func ReadConfig(path string, fs fileSystem) (map[interface{}]interface{}, error) {
	b, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config:=make(map[interface{}]interface{})
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return config, nil
}
