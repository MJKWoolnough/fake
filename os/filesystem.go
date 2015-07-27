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
		},
		make(map[string]FileInfo),
		root,
	}
	cwdmu sync.Mutex
	cwd   *directory
)

func init() {
	Mkdir("/tmp", 0777)
	Cwd("/tmp")
	Chmod("/", 0555)
}

type metadata struct {
	FileMode
	modTime time.Time
	name    string
}

func (m metadata) Name() string {
	return m.name
}

func (m metadata) Mode() FileMode {
	return m.FileMode
}

func (m metadata) ModTime() time.Time {
	return m.modTime
}

func (m *metadata) chmod(fileMode FileMode) {
	m.FileMode = fileMode
}

func (m *metadata) setModTime(m time.Time) {
	m.modTime = m
}

type directory struct {
	metadata
	Contents map[string]FileInfo
	Parent   *directory
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
		},
		make(map[string]FileInfo),
		d,
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

func (d *directory) remove(name string) error {
	if !d.canWrite() {
		return ErrPermissions
	}
	fi, ok := d.Contents[name]
	if !ok {
		return ErrNotExist
	}
	if fi.IsDir() {
		if len(fi.(*directory).Contents) > 0 {
			return ErrNotEmpty
		}
	}
	delete(d.Contents, name)
	return nil
}

func (d *directory) Sys() interface{} {
	return d.Contents
}

type file struct {
	metadata
	Contents []byte
}

func newFile(name string, modTime time.Time, fileMode FileMode, contents []byte) *file {
	return &file{
		metadata{
			name:     name,
			modeTime: modeTime,
			FileMode: fileMode | ModeDir,
			"",
		},
		contents,
	}
}

func (f *file) Sys() interface{} {
	return f.Contents
}
