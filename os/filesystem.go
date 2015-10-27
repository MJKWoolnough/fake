package os

import (
	"io"
	"os"
	"sync"
	"time"
)

var fs = struct {
	sync.Mutex
	root directory
	cwd  path
}{
	root: dir{},
	cwd:  path{},
}

type data interface {
	Size() int64
	Mode() os.FileMode
	ModTime() time.Time
	Data()
}

type node struct {
	name string
	data
}

func (n *node) Name() string {
	return n.name
}

func (n *node) IsDir() bool {
	return n.Mode().IsDir()
}

func (n *node) IsSymlink() bool {
	return n.Mode()&ModeSymlink > 0
}

func (n *node) Sys() interface{} {
	return n.data
}

type file interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	io.ReaderAt
	io.WriterAt
	io.ReaderFrom
	io.WriterTo
	Sync() error
	Truncate(int64) error
	WriteString(string) (int, error)
}

type directory interface {
	io.Seeker
	io.Closer
	Readdir(int) ([]FileInfo, error)
	Readdirnames(int) ([]string, error)
}

type link interface {
	Follow() node
}

type path struct {
	name   string
	deep   uint
	parent *path
}

type modeTime struct {
	FileMode
	modTime time.Time
}

func (m modeTime) Mode() os.FileMode {
	return m.FileMode
}

func (m modeTime) ModTime() time.Time {
	return m.modTime
}

type dir struct {
	modeTime
	contents []FileInfo
}

type fileBytes struct {
	modeTime
	data []byte
}

type fileString struct {
	modeTime
	data string
}
