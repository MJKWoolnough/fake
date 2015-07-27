package os

import "time"

type FileInfo interface {
	Name() string
	Size() int64
	Mode() FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() interface{}
}

func Lstat(name string) (FileInfo, error) {
	return getFile(name)
}

func Stat(name string) (FileInfo, error) {
	return Lstat(name)
}
