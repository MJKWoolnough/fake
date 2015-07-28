package os

import (
	"io"
	"path"
	"sort"
	"unsafe"

	"github.com/MJKWoolnough/memio"
)

const (
	O_RDONLY int = 0x0
	O_WRONLY int = 0x1
	O_RDWR   int = 0x2
	O_APPEND int = 0x400
	O_CREATE int = 0x40
	O_EXCL   int = 0x80
	O_SYNC   int = 0x101000
	O_TRUNC  int = 0x200
)

const (
	SEEK_CUR = 0
	SEEK_SET = 1
	SEEK_END = 2
)

type readWrite struct {
	*memio.ReadWriteMem
}

func (readWrite) Readdir(_ int) ([]FileInfo, error) {
	return nil, ErrInvalid
}
func (readWrite) Readdirnames(_ int) ([]string, error) {
	return nil, ErrInvalid
}

type noWrite struct {
	readWrite
}

func (noWrite) Write(_ []byte) (int, error) {
	return 0, ErrPermission
}

func (noWrite) WriteAt(_ []byte, _ int64) (int, error) {
	return 0, ErrPermission
}

type noRead struct {
	readWrite
}

func (noRead) Read(_ []byte) (int, error) {
	return 0, ErrPermission
}

func (noRead) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, ErrPermission
}

type directoryC struct {
	contents []FileInfo
}

func (d directoryC) Len() int {
	return len(d.contents)
}

func (d directoryC) Less(i, j int) bool {
	return d.contents[i].ModTime().Before(d.contents[i].ModTime())
}

func (d directoryC) Swap(i, j int) {
	d.contents[i], d.contents[j] = d.contents[j], d.contents[i]
}

func (d *directoryC) Readdir(n int) ([]FileInfo, error) {
	if len(d.contents) == 0 {
		return d.contents, io.EOF
	}
	if n > len(d.contents) || n <= 0 {
		c := d.contents
		d.contents = d.contents[len(d.contents):]
		return c, nil
	}
	c := d.contents[:n]
	d.contents = d.contents[n:]
	return c, nil
}

func (d *directoryC) Readdirnames(n int) ([]string, error) {
	dirs, err := d.Readdir(n)
	if err != nil {
		return []string{}, err
	}
	names := make([]string, len(dirs))
	for n, dir := range dirs {
		names[n] = dir.Name()
	}
	return names, nil
}

func (directoryC) Write(_ []byte) (int, error) {
	return 0, ErrInvalid
}

func (directoryC) WriteAt(_ []byte, _ int64) (int, error) {
	return 0, ErrInvalid
}

func (directoryC) Read(_ []byte) (int, error) {
	return 0, ErrInvalid
}

func (directoryC) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, ErrInvalid
}

func (directoryC) Seek(_ int64, _ int) (int64, error) {
	return 0, ErrInvalid
}

type contents interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	Readdir(int) ([]FileInfo, error)
	Readdirnames(int) ([]string, error)
	Seek(int64, int) (int64, error)
	Write([]byte) (int, error)
	WriteAt([]byte, int64) (int, error)
}

type File struct {
	fi   FileInfo
	name string
	contents
}

func Create(name string) (*File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
}

func NewFile(fd uintptr, name string) *File {
	if int(fd) < 0 {
		return nil
	}
	return (*File)(unsafe.Pointer(fd))
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	dir, file := path.Split(path.Clean(name))
	d, err := navigateTo(dir)
	var f FileInfo
	if err == nil {
		f, err = d.get(file)
		if flag&O_CREATE != 0 {
			if IsNotExist(err) {
				f, err = d.create(file, perm)
			} else if flag&O_EXCL != 0 {
				err = ErrExist
			}
		}
	}
	if err != nil {
		return nil, &PathError{
			"open",
			name,
			err,
		}
	}
	var c contents
	if f.IsDir() {
		m := f.Sys().(map[string]FileInfo)
		list := make([]FileInfo, 0, len(m))
		for _, fi := range m {
			list = append(list, fi)
		}
		d := directoryC{list}
		sort.Sort(d)
		c = &d
	} else {
		rw := readWrite{memio.OpenMem(f.Sys().(*[]byte))}
		if flag&O_TRUNC != 0 {
			//rw.Truncate()
		}
		if flag&O_APPEND != 0 {
			rw.Seek(2, 0)
		}
		if flag&O_RDWR != 0 {
			c = rw
		} else if flag&O_WRONLY != 0 {
			c = noRead{rw}
		} else {
			c = noWrite{rw}
		}
	}
	return &File{
		f,
		name,
		c,
	}, nil
}

func Pipe() (*File, *File, error) {
	return nil, nil, ErrUnsupported
}

func (f *File) valid() error {
	if f == nil {
		return ErrInvalid
	}
	if f.fi == nil {
		return ErrClosed
	}
	return nil
}

func (f *File) validPath(op string) error {
	if f == nil {
		return ErrInvalid
	}
	if f.fi == nil {
		return &PathError{
			op,
			f.name,
			ErrClosed,
		}
	}
	return nil
}

func (f *File) Chdir() error {
	if err := f.validPath("chdir"); err != nil {
		return err
	}
	if !f.fi.IsDir() {
		return ErrInvalid
	}
	cwdmu.Lock()
	defer cwdmu.Unlock()
	cwd = f.fi.(*directory)
	return nil
}

func (f *File) Chmod(mode FileMode) error {
	if f == nil {
		return ErrInvalid
	}
	if f.fi == nil {
		return &PathError{
			"chmod",
			f.name,
			ErrClosed,
		}
	}
	type i interface {
		chmod(FileMode)
	}
	f.fi.(i).chmod(mode)
	return nil
}

func (f *File) Chown(_, _ int) error {
	if err := f.validPath("chown"); err != nil {
		return err
	}
	return &PathError{
		"chown",
		f.name,
		ErrUnsupported,
	}
}

func (f *File) Close() error {
	f.fi = nil
	return nil
}

func (f *File) Fd() uintptr {
	if f == nil {
		return ^(uintptr(0))
	}
	return uintptr(unsafe.Pointer(f))
}

func (f *File) Name() string {
	if f == nil {
		return ""
	}
	return f.name
}

func (f *File) Read(b []byte) (int, error) {
	if err := f.validPath("read"); err != nil {
		return 0, err
	}
	return f.contents.Read(b)
}

func (f *File) ReadAt(b []byte, off int64) (int, error) {
	if err := f.validPath("read"); err != nil {
		return 0, err
	}
	return f.contents.ReadAt(b, off)
}

func (f *File) Readdir(n int) ([]FileInfo, error) {
	if err := f.valid(); err != nil {
		return []FileInfo{}, err
	}
	return f.contents.Readdir(n)
}

func (f *File) Readdirnames(n int) ([]string, error) {
	if err := f.valid(); err != nil {
		return []string{}, err
	}
	return f.contents.Readdirnames(n)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if err := f.validPath("seek"); err != nil {
		return 0, err
	}
	return f.contents.Seek(offset, whence)
}

func (f *File) Stat() (FileInfo, error) {
	if err := f.validPath("stat"); err != nil {
		return nil, err
	}
	return f.fi, nil
}

func (f *File) Sync() error {
	if err := f.validPath("fsync"); err != nil {
		return err
	}
	return nil
}

func (f *File) Truncate(size int64) error {
	if err := f.validPath("truncate"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	f.fi.(*file).Contents = nil
	return nil
}

func (f *File) Write(b []byte) (int, error) {
	if err := f.validPath("write"); err != nil {
		return 0, err
	}
	return f.contents.Write(b)
}

func (f *File) WriteAt(b []byte, off int64) (int, error) {
	if err := f.validPath("write"); err != nil {
		return 0, err
	}
	return f.contents.WriteAt(b, off)
}

func (f *File) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}