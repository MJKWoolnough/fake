package os

import (
	"io"
	"os"
	"path"
	"time"
)

func Create(name string) (*File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
}

func NewFile(fd uintptr, name string) *File {
	return nil
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func OpenFile(name string, flag int, perm os.FileMode) (*File, error) {
	n, err := fs.getNode(name, true)
	if err != nil {
		if IsNotExist(err) && flag&O_CREATE != 0 {
			dname, fname := path.Split(name)
			d, err := fs.getDirectory(dname)
			if err != nil {
				return &PathError{
					"open",
					name,
					err,
				}
			}
			n = &fileBytes{
				modeTime: modeTime{
					FileMode(perm),
					time.Now(),
				},
				data: make([]byte, 0),
			}
			if err = d.set(fname, n); err != nil {
				return &PathError{
					"open",
					name,
					err,
				}
			}
		} else {
			return &PathError{
				"open",
				name,
				err,
			}
		}
	} else if flag&O_EXCL == 0 {
		return &PathError{
			"open",
			name,
			ErrExist,
		}
	}

	return nil, nil
}

type File struct {
	name string
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

type readOnly struct {
	file
}

func (readOnly) ReadFrom(io.Reader) (int64, error) {
	return 0, ErrPerm
}

func (readOnly) Sync() error {
	return nil
}

func (readOnly) Truncate(int64) error {
	return ErrPerm
}

func (readOnly) Write([]byte) (int, error) {
	return 0, ErrPerm
}

func (readOnly) WriteAt([]byte, int64) (int, error) {
	return 0, ErrPerm
}

func (readOnly) WriteString(string) (int, error) {
	return 0, ErrPerm
}

type writeOnly struct {
	file
}

func (writeOnly) Read([]byte) (int, error) {
	return 0, ErrPerm
}

func (writeOnly) ReadAt([]byte, int64) (int, error) {
	return 0, ErrPerm
}

func (writeOnly) WriteTo(io.Writer) (int, error) {
	return 0, ErrPerm
}
