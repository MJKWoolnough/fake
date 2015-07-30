package http

import (
	h "net/http"
	oos "os"
	"path"

	"github.com/MJKWoolnough/fake/os"
)

type fileInfo struct {
	os.FileInfo
}

func (f fileInfo) Mode() oos.FileMode {
	return oos.FileMode(f.FileInfo.Mode())
}

type file struct {
	*os.File
}

func (f file) Readdir(n int) ([]oos.FileInfo, error) {
	fis, err := f.File.Readdir(n)
	ofis := make([]oos.FileInfo, len(fis))
	for i, fi := range fis {
		ofis[i] = fileInfo{fi}
	}
	return ofis, err
}

func (f file) Stat() (oos.FileInfo, error) {
	fi, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return fileInfo{fi}, nil
}

type Dir string

func (d Dir) Open(name string) (h.File, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	f, err := os.Open(path.Join(dir, name))
	if err != nil {
		return nil, err
	}
	return file{f}, nil
}
