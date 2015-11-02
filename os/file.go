package os

import (
	"io"
	"os"
	"path"
	"time"
	"unsafe"
)

func Create(name string) (*File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
}

func NewFile(fd uintptr, name string) *File {
	return &File{
		fNode: (*fNode)(unsafe.Pointer(fd)),
		name:  name,
	}
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

type data interface {
	Data() file
}

func OpenFile(name string, flag int, perm os.FileMode) (*File, error) {
	n, err := fs.getNode(name, true)
	if err != nil {
		if IsNotExist(err) && flag&O_CREATE != 0 {
			dname, fname := path.Split(name)
			d, err := fs.getDirectory(dname)
			if err != nil {
				return nil, &PathError{
					"open",
					name,
					err,
				}
			}
			n = bNode{
				&fileBytes{
					modeTime: modeTime{
						FileMode(perm),
						time.Now(),
					},
					data: make([]byte, 0),
				},
				d,
			}
			if err = d.set(fname, n); err != nil {
				return nil, &PathError{
					"open",
					name,
					err,
				}
			}
		} else {
			return nil, &PathError{
				"open",
				name,
				err,
			}
		}
	} else if flag&O_EXCL == 0 {
		return nil, &PathError{
			"open",
			name,
			ErrExist,
		}
	}
	fm := n.Mode()
	switch d := n.node.(type) {
	case *directory:
		if flag != O_RDONLY {
			return nil, &PathError{
				"open",
				name,
				ErrIsDir,
			}
		}
		if fm&0555 == 0 {
			return nil, &PathError{
				"open",
				name,
				ErrPermission,
			}
		}
		contents := make([]dNode, 0, len(d.contents))
		for dname, dnode := range d.contents {
			contents = append(contents, dNode{
				name: dname,
				node: dnode,
			})
		}
		return &File{
			&fNode{
				&dirWrapper{contents: contents},
				n,
			},
			name,
		}, nil
	default:
		var df dfile
		if flag&O_WRONLY != 0 {
			if fm&0222 != 0 {
				df = fileWrapper{writeOnly{n.node.(data).Data()}}
			}
		} else if flag&O_RDWR != 0 {
			if fm&0222 != 0 && fm&0444 != 0 {
				df = fileWrapper{n.node.(data).Data()}
			}
		} else {
			if fm&0444 != 0 {
				df = fileWrapper{readOnly{n.node.(data).Data()}}
			}
		}
		if df == nil {
			return nil, &PathError{
				"open",
				name,
				ErrPermission,
			}
		}
		return &File{
			&fNode{
				df,
				n,
			},
			name,
		}, nil
	}
}

type dfile interface {
	file
	Readdir(int) ([]os.FileInfo, error)
	Readdirnames(int) ([]string, error)
}

type fNode struct {
	dfile
	node bNode
}

type File struct {
	*fNode
	name string
}

func (f File) Chdir() error {
	fi := f.fNode
	if fi == nil {
		return ErrInvalid
	}
	switch d := fi.node.node.(type) {
	case *directory:
		_, fName := path.Split(f.name)
		c := &breadcrumbs{
			name:      fName,
			depth:     f.node.parent.depth + 1,
			previous:  f.node.parent,
			parent:    f.node.parent,
			directory: d,
		}
		fs.Lock()
		fs.cwd = c
		fs.Unlock()
	default:
		return &PathError{
			"chdir",
			f.name,
			ErrIsNotDir,
		}
	}
	return nil
}

func (f File) Chmod(mode os.FileMode) error {
	fi := f.fNode
	if fi == nil {
		return ErrInvalid
	}
	err := fi.node.SetMode(FileMode(mode))
	if err != nil {
		return &PathError{
			"chmod",
			f.name,
			err,
		}
	}
	return nil
}

func (f File) Chown(int, int) error {
	return &PathError{
		"chown",
		f.name,
		ErrUnsupported,
	}
}

func (f *File) Close() error {
	if f == nil {
		return ErrInvalid
	}
	err := f.fNode.Close()
	f.fNode = nil
	return err
}

func (f File) Fd() uintptr {
	return uintptr(unsafe.Pointer(f.fNode))
}

func (f File) Name() string {
	return f.name
}

func (f File) Read(b []byte) (int, error) {
	if f.fNode == nil {
		return 0, ErrInvalid
	}
	return f.fNode.Read(b)
}

func (f File) ReadAt(b []byte, off int64) (int, error) {
	if f.fNode == nil {
		return 0, ErrInvalid
	}
	return f.fNode.ReadAt(b, off)
}

func (f File) Readdir(n int) ([]os.FileInfo, error) {
	if f.fNode == nil {
		return nil, ErrInvalid
	}
	return f.fNode.Readdir(n)
}

func (f File) Readdirnames(n int) ([]string, error) {
	if f.fNode == nil {
		return nil, ErrInvalid
	}
	return f.fNode.Readdirnames(n)
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	if f.fNode == nil {
		return 0, ErrInvalid
	}
	return f.fNode.Seek(offset, whence)
}

func (f File) Stat() (os.FileInfo, error) {
	if f.fNode == nil {
		return nil, ErrInvalid
	}
	return &fileInfo{
		name:    path.Base(f.name),
		size:    f.fNode.node.Size(),
		mode:    f.fNode.node.Mode(),
		modTime: f.fNode.node.ModTime(),
		sys:     f.fNode,
	}, nil
}

func (f File) Sync() error {
	if f.fNode == nil {
		return ErrInvalid
	}
	return f.fNode.Sync()
}

func (f File) Truncate(size int64) error {
	if f.fNode == nil {
		return ErrInvalid
	}
	return f.fNode.Truncate(size)
}

func (f File) Write(b []byte) (int, error) {
	if f.fNode == nil {
		return 0, ErrInvalid
	}
	return f.fNode.Write(b)
}

func (f File) WriteAt(b []byte, offset int64) (int, error) {
	if f.fNode == nil {
		return 0, ErrInvalid
	}
	return f.fNode.WriteAt(b, offset)
}

func (f File) WriteString(s string) (int, error) {
	if f.fNode == nil {
		return 0, ErrInvalid
	}
	return f.fNode.WriteString(s)
}

type fileWrapper struct {
	file
}

func (fileWrapper) Readdir(int) ([]os.FileInfo, error) {
	return nil, ErrIsNotDir
}

func (fileWrapper) Readdirnames(n int) ([]string, error) {
	return nil, ErrIsNotDir
}

type dNode struct {
	name string
	node
}

func (d dNode) Name() string {
	return d.name
}

func (d dNode) IsDir() bool {
	return d.node.Mode().IsDir()
}

func (d dNode) Mode() os.FileMode {
	return os.FileMode(d.node.Mode())
}

func (d dNode) Sys() interface{} {
	return d.node
}

type dirWrapper struct {
	pos      int
	contents []dNode
}

func (dirWrapper) Close() error {
	return nil
}

func (d *dirWrapper) Seek(offset int64, whence int) (int64, error) {
	p := d.pos
	switch whence {
	case SEEK_SET:
		p = whence
	case SEEK_CUR:
		p += whence
	case SEEK_END:
		p = len(d.contents) + whence
	default:
		return 0, ErrInvalid
	}
	if p != 0 {
		return 0, ErrIsDir
	}
	d.pos = 0
	return 0, nil
}

func (d *dirWrapper) Readdir(n int) ([]os.FileInfo, error) {
	if n <= 0 || d.pos+n > len(d.contents) {
		n = len(d.contents) - d.pos
	}
	toRet := make([]os.FileInfo, d.pos+n)
	for i := 0; i < n; i++ {
		toRet[i] = d.contents[d.pos+i]
	}
	d.pos += n
	if len(toRet) == 0 {
		return nil, io.EOF
	}
	return toRet, nil
}

func (d *dirWrapper) Readdirnames(n int) ([]string, error) {
	if n <= 0 || d.pos+n > len(d.contents) {
		n = len(d.contents) - d.pos
	}
	toRet := make([]string, d.pos+n)
	for i := 0; i < n; i++ {
		toRet[i] = d.contents[d.pos+i].name
	}
	d.pos += n
	if len(toRet) == 0 {
		return nil, io.EOF
	}
	return toRet, nil
}

func (dirWrapper) Read([]byte) (int, error) {
	return 0, ErrIsDir
}

func (dirWrapper) ReadAt([]byte, int64) (int, error) {
	return 0, ErrIsDir
}

func (dirWrapper) ReadFrom(io.Reader) (int64, error) {
	return 0, ErrIsDir
}

func (dirWrapper) Sync() error {
	return ErrIsDir
}

func (dirWrapper) Truncate(int64) error {
	return ErrIsDir
}

func (dirWrapper) Write([]byte) (int, error) {
	return 0, ErrIsDir
}

func (dirWrapper) WriteAt([]byte, int64) (int, error) {
	return 0, ErrIsDir
}

func (dirWrapper) WriteString(string) (int, error) {
	return 0, ErrIsDir
}

func (dirWrapper) WriteTo(io.Writer) (int64, error) {
	return 0, ErrIsDir
}

type readOnly struct {
	file
}

func (readOnly) ReadFrom(io.Reader) (int64, error) {
	return 0, ErrPermission
}

func (readOnly) Sync() error {
	return nil
}

func (readOnly) Truncate(int64) error {
	return ErrPermission
}

func (readOnly) Write([]byte) (int, error) {
	return 0, ErrPermission
}

func (readOnly) WriteAt([]byte, int64) (int, error) {
	return 0, ErrPermission
}

func (readOnly) WriteString(string) (int, error) {
	return 0, ErrPermission
}

type writeOnly struct {
	file
}

func (writeOnly) Read([]byte) (int, error) {
	return 0, ErrPermission
}

func (writeOnly) ReadAt([]byte, int64) (int, error) {
	return 0, ErrPermission
}

func (writeOnly) WriteTo(io.Writer) (int64, error) {
	return 0, ErrPermission
}
