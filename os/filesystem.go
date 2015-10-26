package os

import (
	"os"
	"sort"
	"sync"
	"time"

	"github.com/MJKWoolnough/memio"
)

var (
	root = &directory{
		node{
			os.ModeDir | 0777,
			time.Now(),
			"",
			nil,
		},
		make(map[string]os.FileInfo),
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

type dir interface {
	file
	get(string) (os.FileInfo, error)
	set(string, os.FileInfo) error
	removeFile(file)
}

type file interface {
}

type node struct {
	os.FileMode
	modTime time.Time
	name    string
	parent  dir
}

func (n node) Name() string {
	return n.name
}

func (n node) Mode() os.FileMode {
	return n.FileMode
}

func (n node) ModTime() time.Time {
	return n.modTime
}

func (n *node) chmod(fileMode os.FileMode) error {
	n.FileMode = fileMode | (n.FileMode & os.ModeDir)
	return nil
}

func (n *node) setModTime(m time.Time) {
	n.modTime = m
}

func (n *node) move(name string, d dir) error {
	if n.parent == nil {
		return ErrInvalid
	}
	if !canRead(n.parent.FileMode) {
		return ErrPermission
	}
	f, err := n.parent.get(name)
	if err != nil {
		return err
	}
	if err := n.parent.removeFile(name); err != nil {
		return err
	}
	if err := d.set(name, f); err != nil {
		return err
	}
	n.parent = d
	return nil
}

type directory struct {
	node
	Contents map[string]os.FileInfo
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

func (d *directory) create(name string, perm os.FileMode) (os.FileInfo, error) {
	if !canWrite(d.FileMode) {
		return nil, ErrPermission
	}
	if f, ok := d.Contents[name]; ok {
		return f, nil
	}
	if err := namecheck(name); err != nil {
		return nil, err
	}
	f := &bfile{
		node{
			perm &^ os.ModeDir,
			time.Now(),
			name,
			d,
		},
		make([]byte, 0),
	}
	d.Contents[name] = f
	return f, nil
}

func (d *directory) mkdir(name string, fileMode os.FileMode) (*directory, error) {
	if !canWrite(d.FileMode) {
		return nil, ErrPermission
	}
	if _, ok := d.Contents[name]; ok {
		return nil, ErrExist
	}
	if err := namecheck(name); err != nil {
		return nil, err
	}
	e := &directory{
		node{
			fileMode | os.ModeDir,
			time.Now(),
			name,
			d,
		},
		make(map[string]os.FileInfo),
	}
	d.Contents[name] = e
	return e, nil
}

func (d *directory) get(name string) (os.FileInfo, error) {
	if !canExecute(d.FileMode) {
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

func (d *directory) set(name string, f os.FileInfo) error {

}

func (d *directory) remove(name string, all bool) error {
	if !canWrite(d.FileMode) {
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
	list := make([]os.FileInfo, 0, len(d.Contents))
	for _, fi := range d.Contents {
		list = append(list, fi)
	}
	dir := &directoryC{list}
	sort.Sort(dir)
	return dir, nil
}

type bfile struct {
	node
	Contents []byte
}

func (f *bfile) Size() int64 {
	return int64(len(f.Contents))
}

func (f *bfile) Sys() interface{} {
	return &f.Contents
}

func (f *bfile) getContents(flag int) (contents, error) {
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
