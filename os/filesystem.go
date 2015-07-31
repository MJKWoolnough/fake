package os

import (
	"sort"
	"sync"
	"time"

	"github.com/MJKWoolnough/memio"
)

var (
	root = &directory{
		node{
			ModeDir | 0777,
			time.Now(),
			"",
			nil,
		},
		make(map[string]FileInfo),
	}
	cwdmu sync.Mutex
	cwd   *directory
)

func init() {
	root.parent = root
	cwd = root
	Mkdir("/tmp", 0777)
	Chdir("/tmp")
}

type node struct {
	FileMode
	modTime time.Time
	name    string
	parent  *directory
}

func (n node) Name() string {
	return n.name
}

func (n node) Mode() FileMode {
	return n.FileMode
}

func (n node) ModTime() time.Time {
	return n.modTime
}

func (n *node) chmod(fileMode FileMode) error {
	if n.FileMode.isSpecial() {
		return ErrPermission
	}
	n.FileMode = fileMode | (n.FileMode & ModeDir)
	return nil
}

func (n *node) setModTime(m time.Time) {
	n.modTime = m
}

func (n *node) move(name string, d *directory) error {
	if n.parent == nil {
		return ErrInvalid
	}
	if !n.parent.canWrite() || !d.canWrite() {
		return ErrPermission
	}
	f, ok := n.parent.Contents[n.name]
	if !ok {
		return ErrInvalid
	}
	delete(n.parent.Contents, n.name)
	n.parent = d
	d.Contents[name] = f
	return nil
}

type directory struct {
	node
	Contents map[string]FileInfo
}

func namecheck(name string) error {
	for _, c := range name {
		switch c {
		case '\x00', '/':
			return ErrInvalid
		}
	}
	return nil
}

func (d *directory) create(name string, perm FileMode) (FileInfo, error) {
	if !d.canWrite() {
		return nil, ErrPermission
	}
	if f, ok := d.Contents[name]; ok {
		return f, nil
	}
	if err := namecheck(name); err != nil {
		return nil, err
	}
	f := &file{
		node{
			perm &^ ModeDir,
			time.Now(),
			name,
			d,
		},
		make([]byte, 0),
	}
	d.Contents[name] = f
	return f, nil
}

func (d *directory) mkdir(name string, fileMode FileMode) (*directory, error) {
	if !d.canWrite() {
		return nil, ErrExist
	}
	if _, ok := d.Contents[name]; ok {
		return nil, ErrPermission
	}
	if err := namecheck(name); err != nil {
		return nil, err
	}
	e := &directory{
		node{
			fileMode | ModeDir,
			time.Now(),
			name,
			d,
		},
		make(map[string]FileInfo),
	}
	d.Contents[name] = e
	return e, nil
}

func (d *directory) get(name string) (FileInfo, error) {
	if !d.canExecute() {
		return nil, ErrPermission
	}
	switch name {
	case ".":
		return d, nil
	case "..":
		return d.parent, nil
	}
	fi, ok := d.Contents[name]
	if !ok {
		return nil, ErrNotExist
	}
	return fi, nil
}

func (d *directory) remove(name string, all bool) error {
	if !d.canWrite() {
		return ErrPermission
	}
	fi, ok := d.Contents[name]
	if !ok {
		return ErrNotExist
	}
	if fi.IsDir() {
		dir := fi.(*directory)
		if len(dir.Contents) > 0 {
			if all {
				for name := range dir.Contents {
					err := dir.remove(name, true)
					if err != nil {
						return err
					}
				}
			} else {
				return ErrNotEmpty
			}
		}
	}
	delete(d.Contents, name)
	return nil
}

func (d *directory) Size() int64 {
	return 0
}

func (d *directory) Sys() interface{} {
	return d.Contents
}

func (d *directory) getContents(flag int) (contents, error) {
	if flag&O_WRONLY != 0 {
		return nil, ErrIsDir
	}
	list := make([]FileInfo, 0, len(d.Contents))
	for _, fi := range d.Contents {
		list = append(list, fi)
	}
	dir := &directoryC{list}
	sort.Sort(dir)
	return dir, nil
}

type file struct {
	node
	Contents []byte
}

func (f *file) Size() int64 {
	return int64(len(f.Contents))
}

func (f *file) Sys() interface{} {
	return &f.Contents
}

func (f *file) getContents(flag int) (contents, error) {
	rw := readWrite{memio.OpenMem(&f.Contents)}
	if flag&O_TRUNC != 0 {
		f.Contents = f.Contents[:0]
	}
	if flag&O_APPEND != 0 {
		rw.Seek(0, 2)
	}
	if flag&O_RDWR != 0 {
		return rw, nil
	}
	if flag&O_WRONLY != 0 {
		return noRead{rw}, nil
	}
	return noWrite{rw}, nil
}
