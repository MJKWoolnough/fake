package os

type File struct {
}

func Create(name string) (*File, error) {
	return nil, nil
}

func NewFile(fd uintptr, name string) *File {
	return nil
}

func Open(name string) (*File, error) {
	return nil, nil
}

func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	return nil, nil
}

func (f *File) Chdir() error {
	return nil
}

func (f *File) Chmod(mode FileMode) error {
	return nil
}

func (f *File) Chown(uid, git int) error {
	return nil
}

func (f *File) Close() error {
	return nil
}

func (f *File) Fd() uintptr {
	return 0
}

func (f *File) Name() string {
	return ""
}

func (f *File) Read(b []byte) (int, error) {
	return 0, nil
}

func (f *File) ReadAt(b []byte, off int64) (int, error) {
	return 0, nil
}

func (f *File) Readdir(n int) ([]FileInfo, error) {
	return nil, nil
}

func (f *File) Readdirnames(n int) ([]string, error) {
	return nil, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *File) Stat() (FileInfo, error) {
	return nil, nil
}

func (f *File) Sync() error {
	return nil
}

func (f *File) Truncate(size int64) error {
	return nil
}

func (f *File) Write(b []byte) (int, error) {
	return 0, nil
}

func (f *File) WriteAt(b []byte, off int64) (int, error) {
	return 0, nil
}

func (f *File) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}
