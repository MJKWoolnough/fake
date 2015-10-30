package os

import (
	"io"
	"sync"
	"time"
)

type filesystem struct {
	sync.RWMutex
	root, cwd *breadcrumbs
}

var fs = filesystem{
	root: &breadcrumbs{
		name:  "/",
		depth: 1,
		directory: &directory{
			modeTime: modeTime{
				FileMode: FileMode(ModeDir) | 0755,
				modTime:  time.Now(),
			},
			contents: make(map[string]node),
		},
	},
}

func init() {
	fs.root.parent = fs.root
	fs.root.previous = fs.root
	fs.cwd = fs.root
	Mkdir("/dev", 0755)
	Mkdir("/tmp", 0755)
	Chdir("/tmp")
}

/*type special struct {
	fileMode FileMode
	modtime time.Time
	data data
}*/

type node interface {
	Size() int64
	Mode() FileMode
	SetMode(FileMode)
	ModTime() time.Time
	SetModTime(time.Time)
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

func (m *modeTime) SetMode(fm FileMode) {
	m.FileMode = fm
}

func (m *modeTime) SetModTime(t time.Time) {
	m.modTime = t
}

type directory struct {
	modeTime
	sync.RWMutex
	contents map[string]node
}

func newDirectory(fm FileMode) *directory {
	return &directory{
		modeTime: modeTime{
			FileMode: FileMode(ModeDir) | fm,
			modTime:  time.Now(),
		},
		contents: make(map[string]node),
	}
}

func (d *directory) get(name string) (node, error) {
	if d.FileMode&0111 == 0 {
		return nil, ErrPermission
	}
	d.RLock()
	defer d.RUnlock()
	if f, ok := d.contents[name]; ok {
		return f, nil
	}
	return nil, ErrNotExist
}

func (d *directory) set(name string, n node) error {
	if d.FileMode&0333 == 0 {
		return ErrPermission
	}
	d.Lock()
	defer d.Unlock()
	if _, ok := d.contents[name]; ok {
		return ErrExist
	}
	d.contents[name] = n
	return nil
}

func (d *directory) remove(name string) error {
	if d.FileMode&0333 == 0 {
		return ErrPermission
	}
	if _, ok := d.contents[name]; !ok {
		return ErrNotExist
	}
	delete(d.contents, name)
	return nil
}

func (directory) Size() int64 {
	return 0
}

type symlink struct {
	modeTime
	link string
}

func newSymlink(link string) *symlink {
	return &symlink{
		modeTime: modeTime{
			FileMode: FileMode(ModeSymlink) | 077,
			modTime:  time.Now(),
		},
		link: link,
	}
}

func (symlink) Size() int64 {
	return 0
}

type fileBytes struct {
	modeTime
	data []byte
}

type fileString struct {
	modeTime
	data string
}
