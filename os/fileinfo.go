package os

import (
	"os"
	"path"
	"time"
)

type FileMode os.FileMode

func (f FileMode) IsDir() bool {
	return os.FileMode(f).IsDir()
}

func (f FileMode) IsRegular() bool {
	return os.FileMode(f).IsRegular()
}

func (f FileMode) Perm() FileMode {
	return FileMode(os.FileMode(f).Perm())
}

func (f FileMode) String() string {
	return os.FileMode(f).String()
}

func canExecute(f FileMode) bool {
	return f&0111 > 0
}

func canWrite(f FileMode) bool {
	return f&0222 > 0
}

func canRead(f FileMode) bool {
	return f&0444 > 0
}

type fileInfo struct {
	name    string
	size    int64
	mode    FileMode
	modTime time.Time
	sys     interface{}
}

func (f *fileInfo) Name() string {
	return f.name
}

func (f *fileInfo) Size() int64 {
	return f.size
}

func (f *fileInfo) Mode() os.FileMode {
	return os.FileMode(f.mode)
}

func (f *fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *fileInfo) IsDir() bool {
	return f.Mode().IsDir()
}

func (f *fileInfo) Sys() interface{} {
	return f.sys
}

func Lstat(name string) (os.FileInfo, error) {
	return stat(name, false)
}

func Stat(name string) (os.FileInfo, error) {
	return stat(name, true)
}

func stat(name string, followSymlink bool) (os.FileInfo, error) {
	f, err := fs.getNode(name, followSymlink)
	if err != nil {
		return nil, err
	}
	return &fileInfo{
		name:    path.Base(name),
		size:    f.Size(),
		mode:    f.Mode(),
		modTime: f.ModTime(),
		sys:     f,
	}, nil
}
