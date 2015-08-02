package os

import "os"

func Lstat(name string) (os.FileInfo, error) {
	return getFile(name)
}

func Stat(name string) (os.FileInfo, error) {
	return Lstat(name)
}
