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
	cwd  breadcrumbs
}{
	root: &dir{
		modTime:  time.Now(),
		contents: make([]FileInfo, 0),
	},
	cwd: breadcrumbs{
		"",
		0,
		nil,
		nil,
	},
}

func init() {
	fs.cwd.dir = fs.root
	//Mkdir("/dev", 0755)
	//Mkdir("/tmp", 0755)
}

/*type special struct {
	fileMode FileMode
	modtime time.Time
	data data
}*/

type data interface {
	Size() int64
	Mode() FileMode
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

type breadcrumbs struct {
	name   string
	deep   uint
	parent *breadcrumbs
	dir    directory
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
