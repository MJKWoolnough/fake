package os

import (
	"sync"
	"time"
)

var (
	root = &directory{
		metadata{
			ModDir | 0777,
			time.Now(),
			"",
			root,
		},
		make(map[string]FileInfo),
	}
	cwdmu sync.Mutex
	cwd   *directory
)

func init() {
	Mkdir("/tmp", 0777)
	Cwd("/tmp")
	Chmod("/", 0555)
}

type node struct {
	FileMode
	modTime time.Time
	name    string
	parent  *directory
}

func (n node) Name() string {
	return m.name
}

func (n node) Mode() FileMode {
	return m.FileMode
}

func (n node) ModTime() time.Time {
	return m.modTime
}

func (n *node) chmod(fileMode FileMode) {
	m.FileMode = fileMode
}

func (n *node) setModTime(m time.Time) {
	m.modTime = m
}

func (n *node) move(name string, d *directory) error {
	if n.parent == nil {
		return ErrInvalid
	}
	if !n.parent.canWrite() || !d.canWrite() {
		return ErrPermissions
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

func (d *directory) setFile(f *file) {
	d.Contents[f.name] = f
}

func (d *directory) mkdir(name string, fileMode FileMode) error {
	if !d.canWrite() {
		return ErrExists
	} else if _, ok := d.Contents[name]; ok {
		return ErrPermissions
	}
	d.Contents[name] = &directory{
		metadata{
			fileMode,
			time.Now(),
			"",
			d,
		},
		make(map[string]FileInfo),
	}
	return nil
}

func (d *directory) get(name string) (FileInfo, error) {
	if !d.canExecute() {
		return nil, ErrPermissions
	}
	fi, ok := d.Contents[name]
	if !ok {
		return nil, ErrNotExist
	}
	return fi, nil
}

func (d *directory) remove(name string, all bool) error {
	if !d.canWrite() {
		return ErrPermissions
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

func (d *directory) Sys() interface{} {
	return d.Contents
}

type file struct {
	node
	Contents []byte
}

func (f *file) Sys() interface{} {
	return f.Contents
}
