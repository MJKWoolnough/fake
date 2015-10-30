package os

import (
	"io"
	"os"
)

func Create(name string) (*File, error) {

}

func NewFile(fd uintptr, name string) *File {

}

func Open(name string) (*File, error) {

}

func OpenFile(name string, flag int, perm os.FileMode) (*File, error) {

}

type File struct {
}

type fileWrapper struct {
	file
}

func (fileWrapper) Chdir() error {
	return ErrIsNotDir
}

func (fileWrapper) Readdir(int) ([]os.FileInfo, error) {
	return nil, ErrIsNotDir
}

func (fileWrapper) Readdirnames(n int) ([]string, error) {
	return nil, ErrIsNotDir
}

type dirWrapper struct {
	dir      *breadcrumbs
	pos      int
	contents []node
}

func (d *dirWrapper) Chdir() error {
	fs.Lock()
	fs.cwd = d.dir
	fs.Unlock()
	return nil
}

func (d *dirWrapper) Readdir(n int) ([]os.FileInfo, error) {
	return nil, nil
}

func (d *dirWrapper) Readdirnames(n int) ([]string, error) {
	return nil, nil
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

func (dirWrapper) WriteTo(io.Writer) (int, error) {
	return 0, ErrIsDir
}
