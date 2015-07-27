package os

import "unsafe"

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

type File struct {
	fi   FileInfo
	name string
	mode int
	pos  int64
}

func Create(name string) (*File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|_O_TRUNC, 0666)
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
	f, err := getFile(name)
	if err != nil {
		return &PathError{
			"open",
			name,
			err,
		}
	}
	return &File{
		f,
		name,
	}
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
	if err := f.validPath("name"); err != nil {
		return err
	}
	return f.name
}

func (f *File) Read(b []byte) (int, error) {
	if err := f.validPath("read"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	return 0, nil
}

func (f *File) ReadAt(b []byte, off int64) (int, error) {
	if err := f.validPath("read"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	return 0, nil
}

func (f *File) Readdir(n int) ([]FileInfo, error) {
	if err := f.valid(); err != nil {
		return err
	}
	if !f.fi.IsDir() {
		return ErrInvalid
	}
	return nil, nil
}

func (f *File) Readdirnames(n int) ([]string, error) {
	if err := f.valid(); err != nil {
		return err
	}
	if !f.fi.IsDir() {
		return ErrInvalid
	}
	return nil, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if err := f.valid("seek"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	return 0, nil
}

func (f *File) Stat() (FileInfo, error) {
	if err := f.valid("stat"); err != nil {
		return err
	}
	return f.fi, nil
}

func (f *File) Sync() error {
	if err := f.valid("fsync"); err != nil {
		return err
	}
	return nil
}

func (f *File) Truncate(size int64) error {
	if err := f.valid("truncate"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	return nil
}

func (f *File) Write(b []byte) (int, error) {
	if err := f.valid("write"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	return 0, nil
}

func (f *File) WriteAt(b []byte, off int64) (int, error) {
	if err := f.valid("write"); err != nil {
		return err
	}
	if f.fi.IsDir() {
		return ErrInvalid
	}
	return 0, nil
}

func (f *File) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}
