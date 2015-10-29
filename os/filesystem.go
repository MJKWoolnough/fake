package os

import (
	"io"
	"sync"
	"time"
)

type filesystem struct {
	sync.RWMutex
	root, cwd *breadcrumbs
	cwdPath   string
}

var fs = filesystem{
	root: &breadcrumbs{
		directory: &directory{
			modTime:  time.Now(),
			contents: make(map[string]data),
		},
	},
}

func init() {
	fs.root.parent = fs.root
	fs.root.previous = fs.root
	fs.cwd = fs.root
	//Mkdir("/dev", 0755)
	//Mkdir("/tmp", 0755)
}

/*type special struct {
	fileMode FileMode
	modtime time.Time
	data data
}*/

type node interface {
	Size() int64
	Mode() FileMode
	ModTime() time.Time
	Data() interface{}
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

type breadcrumbs struct {
	name             string
	depth            uint
	previous, parent *breadcrumbs
	*directory
}

type modeTime struct {
	FileMode
	modTime time.Time
}

func (m modeTime) Mode() FileMode {
	return m.FileMode
}

func (m modeTime) ModTime() time.Time {
	return m.modTime
}

type directory struct {
	modeTime
	contents map[string]node
}

func (d *directory) get(name string) (node, error) {
	if f, ok := d.contents[name]; ok {
		return f, nil
	}
	return nil, ErrNotExist
}

type symlink struct {
	modeTime
	link string
}

type fileBytes struct {
	modeTime
	data []byte
}

type fileString struct {
	modeTime
	data string
}
