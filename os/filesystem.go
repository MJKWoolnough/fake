package os

import (
	"io"
	"strings"
	"sync"
	"time"

	"github.com/MJKWoolnough/memio"
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
	SetMode(FileMode) error
	ModTime() time.Time
	SetModTime(time.Time) error
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

func (m *modeTime) SetMode(fm FileMode) error {
	m.FileMode = fm
	return nil
}

func (m *modeTime) SetModTime(t time.Time) error {
	m.modTime = t
	return nil
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

func (d *directory) remove(name string, all bool) error {
	if d.FileMode&0333 == 0 {
		return ErrPermission
	}
	d.Lock()
	defer d.Unlock()
	n, ok := d.contents[name]
	if !ok {
		return ErrNotExist
	}
	switch d := n.(type) {
	case *directory:
		if len(d.contents) > 0 {
			if all {
				for c := range d.contents {
					err := d.remove(c, true)
					if err != nil {
						return err
					}
				}
			} else {
				return ErrNotEmpty
			}
		}
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

func (symlink) SetMode(FileMode) error {
	return ErrInvalid
}

func (symlink) Size() int64 {
	return 0
}

type fileBytes struct {
	modeTime
	data []byte
}

func (f *fileBytes) Size() int64 {
	return int64(len(f.data))
}

func (f *fileBytes) Data() file {
	return fileBytesData{memio.OpenMem(&f.data)}
}

type fileBytesData struct {
	*memio.ReadWriteMem
}

func (fileBytesData) Sync() error {
	return nil
}

type fileString struct {
	modeTime
	data string
}

func (f *fileString) SetMode(fm FileMode) error {
	if fm&0222 > 0 {
		return ErrInvalid
	}
	return f.modeTime.SetMode(fm)
}

func (f *fileString) Size() int64 {
	return int64(len(f.data))
}

func (f *fileString) Data() file {
	return fileStringData{strings.NewReader(f.data)}
}

type fileStringData struct {
	*strings.Reader
}

func (f fileStringData) Close() error {
	f.Reader = nil
	return nil
}

func (fileStringData) ReadFrom(io.Reader) (int64, error) {
	return 0, ErrPermission
}

func (fileStringData) Sync() error {
	return nil
}

func (fileStringData) Truncate(int64) error {
	return nil
}

func (fileStringData) Write([]byte) (int, error) {
	return 0, ErrPermission
}

func (fileStringData) WriteAt([]byte, int64) (int, error) {
	return 0, ErrPermission
}

func (fileStringData) WriteString(string) (int, error) {
	return 0, ErrPermission
}
